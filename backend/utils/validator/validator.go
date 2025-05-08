package validator

import "regexp"

// 验证器
// 验证器的作用是验证请求数据的有效性
// 编写网络接口时，需要对请求数据进行验证，会大量使用本包

var (
	// 正则表达式。验证邮箱、手机号、颜色的值是否合法
	EmailRX  = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	PhoneRX  = regexp.MustCompile(`^1[3-9]\d{9}$`)
	ColorHex = regexp.MustCompile("^([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

// 验证是否有效。有效则返回 true
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// 检查是否有效。无效则添加错误信息
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// 检查是否有效。无效则添加错误信息
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// 检查是否有效
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for _, v := range permittedValues {
		if value == v {
			return true
		}
	}

	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)
	for _, v := range values {
		if uniqueValues[v] {
			return false
		}
		uniqueValues[v] = true
	}
	return true
}
