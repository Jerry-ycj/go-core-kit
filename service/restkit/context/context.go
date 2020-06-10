package context

import (
	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/library/stringkit"
	"github.com/mizuki1412/go-core-kit/service/sqlkit"
	"github.com/spf13/cast"
	"net/http"
	"reflect"
	"strings"
)

type Context struct {
	Proxy    iris.Context
	Request  *http.Request
	Response context.ResponseWriter
}

// msg per request
func (ctx *Context) Set(key string, val interface{}) {
	ctx.Proxy.Values().Set(key, val)
}

func (ctx *Context) Get(key string) interface{} {
	return ctx.Proxy.Values().Get(key)
}

func (ctx *Context) DBTx() *sqlkit.Dao {
	var dao *sqlkit.Dao
	if ctx.Get("dbtx") == nil {
		dao = sqlkit.NewTX(ctx.SessionGetSchema())
		ctx.Set("dbtx", dao)
	} else {
		dao = ctx.Get("dbtx").(*sqlkit.Dao)
	}
	return dao
}
func (ctx *Context) DBTxExist() bool {
	return ctx.Get("dbtx") != nil
}

// data: query, form, json/xml, param

// bean 指针
func (ctx *Context) BindForm(bean interface{}) {
	//ctx.Proxy.Params().Get("demo")
	switch bean.(type) {
	case *map[string]interface{}:
		// query会和form合并 post时
		allForm := ctx.Proxy.FormValues()
		for k, v := range allForm {
			(*(bean.(*map[string]interface{})))[k] = v[len(v)-1]
		}
	default:
		//err := ctx.Proxy.ReadForm(bean)
		//if err != nil {
		//	panic(exception.New("form解析错误"))
		//}
		ctx.bindStruct(bean)
		// validator
		err := Validator.Struct(bean)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				panic(exception.New(err.Error()))
			}
			for _, err0 := range err.(validator.ValidationErrors) {
				panic(exception.New("validation failed: " + stringkit.LowerFirst(err0.Field()) + ", need " + err0.Tag()))
			}
		}
	}
}

/// bean 指针
/// description:"xxx" default:"" trim:"true"
func (ctx *Context) bindStruct(bean interface{}) {
	rt := reflect.TypeOf(bean).Elem()
	rv := reflect.ValueOf(bean).Elem()
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		fieldV := rv.Field(i)
		typeString := field.Type.String()
		// multipart file
		if typeString == "class.File" {
			file, header, err := ctx.Proxy.FormFile(stringkit.LowerFirst(field.Name))
			if err != nil {
				panic(exception.New(err.Error()))
			}
			fieldV.Set(reflect.ValueOf(class.File{
				File:   file,
				Header: header,
			}))
			continue
		}
		// bind
		val := ctx.Proxy.FormValue(stringkit.LowerFirst(field.Name))
		if val == "" {
			val = ctx.Proxy.URLParam(stringkit.LowerFirst(field.Name))
		}
		// 判断trim
		if field.Tag.Get("trim") == "true" {
			val = strings.TrimSpace(val)
		}
		if val == "" {
			// default
			val = field.Tag.Get("default")
		}
		switch typeString {
		case "string":
			fieldV.SetString(val)
		case "int32", "int", "int64":
			if !stringkit.IsNull(val) {
				fieldV.SetInt(cast.ToInt64(val))
			}
		case "float64":
			if !stringkit.IsNull(val) {
				fieldV.SetFloat(cast.ToFloat64(val))
			}
		case "class.Int32":
			if !stringkit.IsNull(val) {
				tmp := class.Int32{Int32: cast.ToInt32(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Int64":
			if !stringkit.IsNull(val) {
				tmp := class.Int64{Int64: cast.ToInt64(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Float64":
			if !stringkit.IsNull(val) {
				tmp := class.Float64{Float64: cast.ToFloat64(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.Bool":
			if !stringkit.IsNull(val) {
				tmp := class.Bool{Bool: cast.ToBool(val), Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.String":
			if !stringkit.IsNull(val) {
				tmp := class.String{String: val, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.ArrInt":
			if !stringkit.IsNull(val) {
				var p []int32
				jsonkit.ParseObj(val, &p)
				tmp := class.ArrInt{Array: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.ArrString":
			if !stringkit.IsNull(val) {
				var p []string
				jsonkit.ParseObj(val, &p)
				tmp := class.ArrString{Array: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		case "class.MapString":
			if !stringkit.IsNull(val) {
				var p map[string]interface{}
				jsonkit.ParseObj(val, &p)
				tmp := class.MapString{Map: p, Valid: true}
				fieldV.Set(reflect.ValueOf(tmp))
			}
		}
	}
}
