package restkit

import (
	ctx "context"
	"github.com/gin-gonic/gin"
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/init/configkey"
	"github.com/mizuki1412/go-core-kit/service/configkit"
	"github.com/mizuki1412/go-core-kit/service/logkit"
	"github.com/mizuki1412/go-core-kit/service/restkit/context"
	"github.com/mizuki1412/go-core-kit/service/restkit/middleware"
	router2 "github.com/mizuki1412/go-core-kit/service/restkit/router"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var router *router2.Router
var server *http.Server

func defaultEngine() {
	if !configkit.GetBoolD(configkey.ProfileDev) {
		gin.SetMode(gin.ReleaseMode)
	}
	router = &router2.Router{
		Proxy: gin.New(),
	}
	router.Proxy.Use(context.InitSession())
	router.Use(middleware.Log())
	router.Use(middleware.Cors())
	router.Use(middleware.Recover())
	//router.Use(cors.Default())

	if configkit.GetBool(configkey.RestPPROF, false) {
		// todo  p := pprof.New()
	}
	// max request size todo
	//router.Proxy.Use(iris.LimitRequestBodySize(int64(configkit.GetInt(configkey.RestRequestBodySize, 100)) << 20))
	// 其他错误如404，
	//router.OnError(middleware.Cors())
	// add base path
	base := configkit.GetStringD(configkey.RestServerBase)
	if base != "" {
		if base[0] != '/' {
			base = "/" + base
		}
		if base[len(base)-1] == '/' {
			base = base[:len(base)-1]
		}
		router.Base = base
	}
}

func Run() error {
	if router == nil {
		defaultEngine()
	}
	port := configkit.GetString(configkey.RestServerPort, "8080")
	router.RegisterSwagger()
	server = &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	go func() {
		logkit.Info("Listening and serving HTTP on " + port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logkit.Fatal(exception.New(err.Error()))
		}
	}()
	// https://github.com/gin-gonic/examples/blob/master/graceful-shutdown/graceful-shutdown/notify-without-context/server.go
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logkit.Info("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctxt, cancel := ctx.WithTimeout(ctx.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctxt); err != nil {
		logkit.Error(exception.New(err.Error()))
		return err
	}
	return nil
}

func Shutdown() {
	if server != nil {
		err := server.Shutdown(ctx.Background())
		if err != nil {
			logkit.Error(exception.New(err.Error()))
		}
	}
}

// AddActions 导入业务模块，其中的路由和中间件
func AddActions(actionInits ...func(r *router2.Router)) {
	if router == nil {
		defaultEngine()
	}
	for _, action := range actionInits {
		action(router)
	}
}

func GetRouter() *router2.Router {
	if router == nil {
		defaultEngine()
	}
	return router
}
