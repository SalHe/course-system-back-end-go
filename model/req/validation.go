package req

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/se2022-qiaqia/course-system/log"
	"github.com/se2022-qiaqia/course-system/model/resp"
	"reflect"
	"regexp"
)

var translator ut.Translator
var (
	regexUsername  = regexp.MustCompile(`^\w{4,32}$`)
	regexpPassword = regexp.MustCompile(`^\w{6,16}$`)
)

func InitValidation() {
	if validate, ok := binding.Validator.Engine().(*validator.Validate); ok {
		initTranslation(validate)
		initGetTag(validate)
		initCustomValidation(validate)
	}
}

func initCustomValidation(validate *validator.Validate) {
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		if ok, err := regexp.MatchString(`[a-zA-Z]`, fl.Field().String()); !ok || err != nil {
			return false
		}
		return regexUsername.MatchString(fl.Field().String())
	})
	validate.RegisterTranslation("username", translator,
		func(ut ut.Translator) error {
			return ut.Add("username", "{0}必须是4-32位字母或数字组成", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("username", fe.Field())
			return t
		},
	)

	validate.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return regexpPassword.MatchString(fl.Field().String())
	})
	validate.RegisterTranslation("password", translator,
		func(ut ut.Translator) error {
			return ut.Add("password", "{0}必须由6-16位数字或字母组成", false)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("password", fe.Field())
			return t
		},
	)
}

func initGetTag(validate *validator.Validate) {
	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		description := field.Tag.Get("description")
		if len(description) > 0 {
			return fmt.Sprintf("%v(%v)", description, field.Name)
		}
		return field.Name
	})
}

func initTranslation(validate *validator.Validate) {
	zhLocales := zh_Hans.New()
	uni := ut.New(zhLocales)
	translator, _ = uni.GetTranslator("zh")
	err := zh.RegisterDefaultTranslations(validate, translator)
	if err != nil {
		log.L.Fatal().
			Err(err).
			Msg("初始化验证器失败")
		return
	}
	log.L.Info().Msg("初始化验证器成功")
	return
}

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBind(obj); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			var finalMsg string
			for _, msg := range errs.Translate(translator) {
				finalMsg += /*key + ":" +*/ msg + ";\n"
			}
			resp.FailJust(finalMsg, c)
			return false
		}
	}
	return true
}
