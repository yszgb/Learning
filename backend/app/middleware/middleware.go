package middleware

import (
	"expvar"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"learning.com/learning/models"
	"learning.com/learning/utils/helper"
	"learning.com/learning/utils/validator"
)

type Middleware struct {
	helper      *helper.Helper
	limitConfig *models.LimiterConfig
	models      *models.Models
}

func New(helper *helper.Helper, limitConfig *models.LimiterConfig, models *models.Models) *Middleware {
	return &Middleware{
		helper:      helper,
		limitConfig: limitConfig,
		models:      models,
	}
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		// 如果没有提供 Authorization 头，将用户设置为匿名用户并继续处理请求
		if authorizationHeader == "" {
			r = m.ContextSetUser(r, models.AnnoymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.helper.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		v := validator.New()
		if models.ValidateTokenPlaintext(v, token); !v.Valid() {
			m.helper.InvalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := m.models.Users.GetUserForToken(models.ScopeAuthentication, token)
		if err != nil {
			switch {
			case err == helper.ErrRecordNotFound:
				m.helper.InvalidAuthenticationTokenResponse(w, r)
			default:
				m.helper.ServerErrorResponse(w, r, err)
			}
			return
		}
		r = m.ContextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

// 判断匿名用户登录中间件
func (m *Middleware) RequirAuthenticateUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := m.ContextGetUser(r)
		if user.IsAnonymous() {
			// 如果用户未登录，返回错误
			m.helper.AuthenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// 记录访问情况、运行情况
func (m *Middleware) Metrics(next http.Handler) http.Handler {
	totoalRequestsReceived := expvar.NewInt("total_requests_received")
	totoalResponsesSent := expvar.NewInt("total_responses_sent")
	totoalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	totoalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totoalRequestsReceived.Add(1)
		metrics := httpsnoop.CaptureMetrics(next, w, r)
		totoalResponsesSent.Add(1)
		totoalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totoalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

// 限流中间件
//
// 检查速率限制器是否启用。如果启用，获取客户端的 IP 地址。
//
// 使用互斥锁保护 clients 映射。检查 IP 是否存在于映射中。
// 如果不存在，创建一个新的速率限制器，并将其添加到映射中。如果存在，更新最后访问时间。
//
// 检查速率限制器是否允许请求。如果不允许，返回速率限制错误。如果允许，调用下一个处理程序处理请求。
func (m *Middleware) RateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// 创建协程。每 1 分钟检查一次 clients 中的元素，如果超过 3 分钟没有被访问，就删除
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.limitConfig.Enabled {
			ip := realip.FromRequest(r)
			mu.Lock()
			if _, found := clients[ip]; !found {
				rps := m.limitConfig.Rps
				burst := m.limitConfig.Burst
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(rps), burst),
				}
			}
			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				m.helper.RateLimitExceededResponse(w, r)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		msg := fmt.Sprintf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		m.helper.Logger.Info(msg, nil)
		next.ServeHTTP(w, r)
	})
}

// 捕获 panic，并处理
//
// 1. 包装传入的 http.Handler，返回新的 http.Handler
//
// 2. 新的 http.Handler 先调用 defer 语句捕获 panic
//
// 如果捕获到，设置响应头 Connection 为 close，并调用 ServerErrorResponse 处理错误
//
// 如果没有捕获到，正常调用传入的 http.Handler 的 ServeHTTP 方法处理请求
//
// 提高健壮性，防止未捕获的 panic 导致服务器崩溃
func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				m.helper.ServerErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
