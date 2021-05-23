package api

import (
	"github.com/chris-joseph/golang-ecs/pkg/config"
	"github.com/chris-joseph/golang-ecs/pkg/data"
	"github.com/chris-joseph/golang-ecs/pkg/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

type App struct{
	server *echo.Echo
	userSvc services.IUserService
	cfg     *config.Settings
}

func New(cfg *config.Settings,client *mongo.Client)*App  {
	server:=echo.New()

	server.Use(middleware.Recover())
	server.Use(middleware.RequestID())

	userProvider:=data.NewUserProvider(cfg,client)

	userSvc:=services.NewUserService(cfg,userProvider)


	return &App{
		server: server,
		userSvc: userSvc,
		cfg: cfg,
	}
}

func (a App)ConfigureRoutes()  {
	a.server.GET("/v1/public/healthy",a.HealthCheck)
	a.server.GET("/v1/public/register",a.Register)
	a.server.GET("/v1/public/login",a.Login)
}

func (a App)Start(){
	a.ConfigureRoutes()
	a.server.Start(":5000")
}