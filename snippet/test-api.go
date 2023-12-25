package snippet

import (
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/library/jsonkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/openapi"
	"github.com/mizuki1412/go-core-kit/service/restkit/router"
)

func Init(router *router.Router) {
	tag := "user:用户模块"
	router.Get("/param", test).Api(openapi.Tag(tag), openapi.Summary("test1"), openapi.ReqParam(testParam{}))
	router.Post("/param", test).Api(openapi.Tag(tag), openapi.Summary("test2"), openapi.ReqBody(testParam{}))
	router.Post("/post/file", file).Api(openapi.Tag(tag), openapi.Summary("test-file"), openapi.ReqBody(fileParams{}))
	router.Put("/put", test).Api(openapi.Tag(tag), openapi.Summary("test3"), openapi.ReqParam(testParam{}))
	router.Delete("/delete", test).Api(openapi.Tag(tag), openapi.Summary("test4"), openapi.ReqParam(testParam{}))
}

type testParam struct {
	Id        int32        `comment:"标识" validate:"required"`
	ValStr    class.String `comment:"数值"`
	ValLong   int64
	ValDouble class.Float64
}

func test(ctx *context.Context) {
	params := testParam{}
	ctx.BindForm(&params)
	println(jsonkit.ToString(params))
	ctx.JsonSuccess()
}

type fileParams struct {
	Id   string
	File class.File `validate:"required"`
}

func file(ctx *context.Context) {
	params := fileParams{}
	ctx.BindForm(&params)
	println(len(params.File.GetBytes()))
	ctx.JsonSuccess()
}
