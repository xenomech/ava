package validator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	instance *validator.Validate
	once     sync.Once
)

func get() *validator.Validate {
	once.Do(func() {
		instance = validator.New()

		instance.RegisterTagNameFunc(func(field reflect.StructField) string {
			name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
			switch name {
			case "-":
				return ""
			case "":
				return field.Name
			default:
				return name
			}
		})
	})

	return instance
}

func ValidateStruct(payload any) error {
	return get().Struct(payload)
}

func FirstError(err error) string {
	if err == nil {
		return ""
	}

	var errs validator.ValidationErrors
	if errors.As(err, &errs) && len(errs) > 0 {
		return errs[0].Field() + " " + describe(errs[0])
	}

	return "Invalid request"
}

func FieldErrors(err error) map[string]string {
	var errs validator.ValidationErrors
	if !errors.As(err, &errs) {
		return nil
	}

	fields := make(map[string]string, len(errs))

	for _, e := range errs {
		name := e.Field()
		if name == "" {
			continue
		}

		if _, seen := fields[name]; !seen {
			fields[name] = describe(e)
		}
	}

	return fields
}

var fixedMessages = map[string]string{
	"required": "is required",
	"email":    "must be a valid email address",
	"url":      "must be a valid URL",
	"uuid":     "must be a valid identifier",
	"uuid4":    "must be a valid identifier",
	"alphanum": "must contain only letters and numbers",
	"dive":     "contains an invalid entry",
}

func describe(e validator.FieldError) string {
	if message, ok := fixedMessages[e.Tag()]; ok {
		return message
	}

	return describeWithParam(e)
}

func describeWithParam(e validator.FieldError) string {
	param := e.Param()

	switch e.Tag() {
	case "min":
		return "must be at least " + param + lengthUnit(e)
	case "max":
		return "must be at most " + param + lengthUnit(e)
	case "len":
		return "must be exactly " + param + " characters"
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(param, " ", ", ")
	case "eqfield":
		return "must match " + param
	case "gte":
		return "must be " + param + " or more"
	case "lte":
		return "must be " + param + " or less"
	}

	if _, err := strconv.Atoi(param); param != "" && err == nil {
		return "is invalid (" + e.Tag() + " " + param + ")"
	}

	return "is invalid"
}

func lengthUnit(e validator.FieldError) string {
	if e.Kind() == reflect.String {
		return " characters"
	}

	return ""
}
