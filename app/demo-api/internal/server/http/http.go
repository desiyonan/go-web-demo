package http

import (
	"fmt"
	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	_ "go-web-demo/app/demo-api/docs"
	"go-web-demo/app/demo-api/internal/conf"
	"go-web-demo/app/demo-api/internal/middleware/access"
	"go-web-demo/app/demo-api/internal/service"
	"go-web-demo/library/log"
	"net/http"
	"time"
)

var (
	srv *service.Service
)

// New init
func New(c *conf.Config, s *service.Service) (httpSrv *http.Server) {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(access.SlowAccess())
	route(engine)
	readTimeout := conf.Conf.App.ReadTimeout
	writeTimeout := conf.Conf.App.WriteTimeout
	endPoint := fmt.Sprintf(":%d", conf.Conf.App.HttpPort)
	maxHeaderBytes := 1 << 20
	httpSrv = &http.Server{
		Addr:           endPoint,
		Handler:        engine,
		ReadTimeout:    time.Duration(readTimeout),
		WriteTimeout:   time.Duration(writeTimeout),
		MaxHeaderBytes: maxHeaderBytes,
	}
	srv = s
	ginpprof.Wrapper(engine)
	go func() {
		// service connections
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("srv.ListenAndServe() error(%v) | config(%v)", err, c)
			panic(err)
		}
	}()
	return
}

func route(e *gin.Engine) {
	e.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	user := e.Group("/api/v1/user")
	{
		//获取当前用户
		user.GET("/current", CurrentUser)
		//获取用户假期
		user.GET("/leaves", CurrentUser)
		//更新用户
		user.PUT("/update", UpdateUser)
	}

	leaveProcess := e.Group("/api/v1/leave-process")
	{

		//获取用户列表
		leaveProcess.GET("/last", LastLeaveProcess)
		//获取用户列表
		leaveProcess.GET("/list", ListLeaveProcess)
		// 添加用户
		leaveProcess.POST("/add", AddLeaveProcess)
		// 添加用户
		leaveProcess.POST("/detail", GetLeaveProcess)
	}
}
