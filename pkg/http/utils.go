package http

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func BindData(ctx *gin.Context, obj any) error {
	if ctx.Request.ContentLength > 0 {
		if err := ctx.ShouldBindJSON(obj); err != nil {
			return err
		}

		SanitizeStruct(obj)
		return nil
	}

	if err := ctx.ShouldBind(obj); err != nil {
		return err
	}

	return nil
}

func BindParams[T any](ctx *gin.Context, key string) (T, error) {
	var res T
	var value string

	if key == "ids" {
		value = ctx.Query(key)
		if value == "" {
			return res, errors.New("ids are required")
		}
		res = any(strings.Split(value, ",")).(T)
		return res, nil
	}

	value = ctx.Param(key)
	if value == "" {
		return res, fmt.Errorf("%s is required", key)
	}

	switch any(res).(type) {
	case string:
		res = any(value).(T)
	case int64:
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return res, err
		}
		res = any(id).(T)
	default:
		return res, fmt.Errorf("unsupported type")
	}

	return res, nil
}

func SanitizeStruct(payload interface{}) {
	v := reflect.ValueOf(payload).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			sanitizedValue := strings.TrimSpace(field.String())
			field.SetString(sanitizedValue)
		}
	}
}
