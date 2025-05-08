package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
	ErrDuplicateEmail = errors.New("duplicate email")
)

func (helper *Helper) LogError(r *http.Request, err error) {
	helper.Logger.Error(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (helper *Helper) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	data := map[string]interface{}{
		"message": message,
	}
	resp := helper.NewResponse(-1, data)
	err := helper.WriteJSON(w, status, resp, nil)
	if err != nil {
		helper.LogError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (helper *Helper) ForbiddenResponse(w http.ResponseWriter, r *http.Request) {
	message := "You are not authorized to access this resource"
	helper.ErrorResponse(w, r, http.StatusForbidden, message)
}

func (helper *Helper) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	helper.LogError(r, err)
	message := "The server encountered a problem and could not process your request"
	helper.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (helper *Helper) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	helper.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (helper *Helper) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := "The requested method is not allowed for this resource"
	helper.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (helper *Helper) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	helper.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (helper *Helper) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	var strBuilder strings.Builder // 用于高效构建字符串
	// 适合循环、多次使用 Sprintf()
	// '+' 拼接字符串时会创建新的字符串对象，导致内存分配和复制，因此不适合频繁拼接字符串
	for key, value := range errors {
		strBuilder.WriteString(fmt.Sprintf("%s: %s, ", key, value))
	}
	result := strBuilder.String()
	if len(result) > 0 {
		result = result[:len(result)-2] // 去掉最后一个逗号和空格
		result += "."
	}
	helper.ErrorResponse(w, r, http.StatusUnprocessableEntity, result)
}

func (helper *Helper) EditConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "Unable to update the record due to an edit conflict, please try again"
	helper.ErrorResponse(w, r, http.StatusConflict, message)
}

func (helper *Helper) RateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "Rate limit exceeded"
	helper.ErrorResponse(w, r, http.StatusTooManyRequests, message)
}

func (helper *Helper) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid authentication credentials"
	helper.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (helper *Helper) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	message := "Invalid or missing authentication token"
	helper.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (helper *Helper) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "You must be authenticated to access this resource"
	helper.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (helper *Helper) NotPermittedResponse(w http.ResponseWriter, r *http.Request) {
	message := "Your account does not have the necessary permissions to access this resource"
	helper.ErrorResponse(w, r, http.StatusForbidden, message)
}
