package context

import (
	"github.com/gin-gonic/gin/render"
	"github.com/mizuki1412/go-core-kit/class"
	"github.com/mizuki1412/go-core-kit/service/storagekit"
	"net/http"
	"net/url"
)

type RestRet struct {
	Result  int          `json:"result"`
	Message class.String `json:"message,omitempty"`
	Data    any          `json:"data,omitempty"`
	// 分页信息
	Total class.Int64 `json:"total,omitempty" description:"总数，如果data是列表并且分页"`
}

const ResultErr = 0
const ResultSuccess = 1
const ResultAuthErr = 2
const ResultUnauthorized = 403

// TransferRestRet 用于自定义返回结构时的转换
var TransferRestRet = func(ret RestRet) any {
	return ret
}

// Json http返回json数据
func (ctx *Context) Json(ret RestRet) {
	var code int
	switch ret.Result {
	case ResultSuccess:
		code = http.StatusOK
	case ResultAuthErr, ResultUnauthorized:
		code = http.StatusUnauthorized
	default:
		code = http.StatusBadRequest
	}
	ctx.Proxy.JSON(code, TransferRestRet(ret))
}

func (ctx *Context) JsonSuccess(data any) {
	ctx.Json(RestRet{
		Result: ResultSuccess,
		Data:   data,
	})
}

func (ctx *Context) RawSuccess(data []byte) {
	ctx.Proxy.Render(http.StatusOK, render.Data{Data: data})
}

func (ctx *Context) Html(data []byte) {
	ctx.Proxy.Render(http.StatusOK, render.Data{Data: data, ContentType: "text/html"})
}

// JsonSuccessWithPage 带分页信息
func (ctx *Context) JsonSuccessWithPage(data any, total uint64) {
	ret := RestRet{
		Result: ResultSuccess,
		Data:   data,
	}
	ret.Total.Set(total)
	ctx.Json(ret)
}
func (ctx *Context) JsonError(msg string) {
	ctx.Json(RestRet{
		Result:  ResultErr,
		Message: class.NewString(msg),
	})
}

func (ctx *Context) SetFileHeader(filename string) {
	ctx.Proxy.Header("Content-Disposition", "attachment; filename="+url.QueryEscape(filename))
	ctx.Proxy.Header("Content-Type", "application/octet-stream")
	ctx.Proxy.Header("Content-Transfer-Encoding", "binary")
	ctx.Proxy.Header("Pragma", "No-cache")
	ctx.Proxy.Header("Cache-Control", "No-cache")
	ctx.Proxy.Header("Expires", "0")
}
func (ctx *Context) SetJsonHeader() {
	ctx.Proxy.Header("Content-Type", "application/json")
}
func (ctx *Context) FileRaw(data []byte, name string) {
	ctx.SetFileHeader(name)
	ctx.RawSuccess(data)
}

// File 相对于项目目录路径的
func (ctx *Context) File(relativePath, name string) {
	ctx.FileDirect(storagekit.GetFullPath(relativePath), name)
}
func (ctx *Context) File2(relativePathName string) {
	ctx.Proxy.File(storagekit.GetFullPath(relativePathName))
}

func (ctx *Context) FileDirect(obsolutePath, name string) {
	ctx.Proxy.File(obsolutePath + name)
}
