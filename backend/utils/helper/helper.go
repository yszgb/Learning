package helper

import (
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
	"learning.com/learning/utils/backgroud"
	"learning.com/learning/utils/jsonlog"
	"learning.com/learning/utils/validator"
)

type Helper struct {
	rand       *rand.Rand // 用于生成随机数
	Logger     *jsonlog.Logger
	Background *backgroud.Background // 用于执行后台任务
}

func New(logger *jsonlog.Logger) *Helper {
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &Helper{
		rand:       rand,
		Logger:     logger,
		Background: backgroud.New(logger),
	}
}

// ReadIDParam 从 URL 路径中读取 id 参数
func (helper *Helper) ReadIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	} else {
		return id, nil
	}
}

func (helper *Helper) ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)

	// 对前端请求进行错误处理
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return errors.New("body contains badly-formed JSON (at character " + strconv.Itoa(int(syntaxError.Offset)) + ")")
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return errors.New("body contains incorrect JSON type for field " + unmarshalTypeError.Field)
			}
			return errors.New("body contains incorrect JSON type (at character " + strconv.Itoa(int(unmarshalTypeError.Offset)) + ")")

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return errors.New("body contains unknown key " + fieldName)

		case errors.As(err, &maxBytesError):
			return errors.New("body must not be larger than " + strconv.FormatInt(maxBytesError.Limit, 10) + " bytes")

		case errors.As(err, &invalidUnmarshalError):
			panic(err)
		default:
			return err
		}
	}

	// 判断是否包含多个 JSON 结构
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

type Response map[string]interface{}

// 装配返回给前端的数据
func (helper *Helper) NewResponse(code int, data map[string]interface{}) Response {
	response := Response(data)
	response["code"] = code
	return response
}

func (helper *Helper) WriteJSON(w http.ResponseWriter, status int, data Response, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}
	js = append(js, '\n')
	for key, value := range headers {
		w.Header()[key] = value
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (helper *Helper) ReadString(qs url.Values, key, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

func (helper *Helper) ReadCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs[key]
	if len(csv) == 0 {
		return defaultValue
	}
	return strings.Split(csv[0], ",")
}

func (helper *Helper) ReadInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (helper *Helper) ReadInt64(qs url.Values, key string, defaultValue int64, v *validator.Validator) int64 {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}
	return i
}

func (helper *Helper) ReadBool(qs url.Values, key string, defaultValue bool, v *validator.Validator) bool {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	b, err := strconv.ParseBool(s)
	if err != nil {
		v.AddError(key, "must be a boolean value")
		return defaultValue
	}
	return b
}

func (helper *Helper) ReadInt64Slice(qs url.Values, key string, v *validator.Validator) []int64 {
	parts, ok := qs[key]
	if !ok || len(parts) == 0 {
		return []int64{}
	}
	ints := make([]int64, 0, len(parts))
	for _, part := range parts {
		i, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			v.AddError(key, "must be a list of integers")
			return []int64{}
		}
		ints = append(ints, i)
	}
	return ints
}

// GeneratePasswordHash 生成密码的哈希值
func (helper *Helper) GeneratePasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), 12)
}

func (helper *Helper) ComparePasswordAndHash(password string, hash []byte) (bool, error) {
	err := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	} else {
		return true, nil
	}
}

// RandomString 生成随机字符串
func (helper *Helper) RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[helper.rand.Intn(len(charset))]
	}
	return string(b)
}
