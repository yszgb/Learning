package api

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"learning.com/learning/models"
	"learning.com/learning/utils/validator"
)

// 登录
//
// 前端访问这个接口
func (api *Api) login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Account string `json:"account"`
		Code    string `json:"code"`
	}
	err := api.helper.ReadJSON(w, r, &input)
	if err != nil {
		api.helper.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(validator.Matches(input.Account, validator.EmailRX), "account", "must be a valid email address")
	v.Check(input.Code != "", "code", "must be provided")
	if !v.Valid() {
		api.helper.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := api.models.Users.GetUserBriefByAccountOrID(input.Account, 0)
	if err != nil {
		api.helper.InvalidCredentialsResponse(w, r)
		return
	}

	if user.PasswdHash == nil {
		api.helper.ErrorResponse(w, r, http.StatusUnauthorized, "password not set")
		return
	}

	math, err := api.helper.ComparePasswordAndHash(input.Code, user.PasswdHash)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
		return
	} else if !math {
		api.helper.InvalidCredentialsResponse(w, r)
		return
	}

	// 账号存在，密码正确，生成 token
	// token 存在一个月的有效期
	token, err := api.models.Tokens.New(user.ID, 30*24*time.Hour, models.ScopeAuthentication)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
		return
	}

	data := map[string]interface{}{
		"token": token.Plaintext,
		"brief": user,
	}

	resp := api.helper.NewResponse(0, data)
	err = api.helper.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
	}
}

func (api *Api) logout(w http.ResponseWriter, r *http.Request) {
	authorization := r.Header.Get("Authorization")
	headerParts := strings.Split(authorization, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		api.helper.InvalidCredentialsResponse(w, r)
		return
	}

	token := headerParts[1]
	err := api.models.Tokens.DeleteToken(token)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
		return
	}

	data := map[string]interface{}{}
	resp := api.helper.NewResponse(0, data)
	err = api.helper.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
	}
}

// 向前端返回错误
func (api *Api) getUserBrief(w http.ResponseWriter, r *http.Request) {
	user := api.middleware.ContextGetUser(r)
	brief, err := api.models.Users.GetUserBriefByAccountOrID("", user.ID)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
		return
	}

	data := map[string]interface{}{
		"brief": brief,
	}
	resp := api.helper.NewResponse(0, data)
	err = api.helper.WriteJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		api.helper.ServerErrorResponse(w, r, err)
	}
}

// 头像 接口
func (api *Api) getUserAvatar(w http.ResponseWriter, r *http.Request) {
	qs := r.URL.Query() // 获取查询参数
	avatar := api.helper.ReadString(qs, "k", "")
	if avatar == "" {
		user := api.middleware.ContextGetUser(r)
		avatarFromDB, err := api.models.Users.GetUserAvatarById(user.ID)
		if err != nil {
			api.helper.ServerErrorResponse(w, r, err)
			return
		}
		if avatarFromDB == nil {
			api.helper.NotFoundResponse(w, r)
			return
		}
		avatar = *avatarFromDB
	}
	filepath := filepath.Join(api.config.Path.AvatarsPath(), avatar)
	http.ServeFile(w, r, filepath) // 直接返回文件。前端会自动解析
}

func (api *Api) userRoutes(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/v1/user/login", api.login)
	router.HandlerFunc(http.MethodPost, "/v1/user/logout", api.middleware.RequirAuthenticateUser(api.logout))

	router.HandlerFunc(http.MethodGet, "/v1/user/brief", api.middleware.RequirAuthenticateUser(api.getUserBrief))

	router.HandlerFunc(http.MethodGet, "/v1/user/avatar", api.getUserAvatar)
}
