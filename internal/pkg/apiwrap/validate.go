package apiwrap

import (
	"regexp"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/go-playground/validator/v10"
)

// ValidateBsonId 验证bson.ObjectId
var ValidateBsonId validator.Func = func(fl validator.FieldLevel) bool {
	id := fl.Field().String()
	if _, err := bson.ObjectIDFromHex(id); err != nil {
		return false
	}
	return true
}

// ValidateVersion 验证文档版本号
var ValidateVersion validator.Func = func(fl validator.FieldLevel) bool {
	version := fl.Field().String()
	if matched, _ := regexp.MatchString(`^\d+\.\d+$`, version); !matched {
		return false
	}
	return true
}
