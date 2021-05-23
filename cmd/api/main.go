package main

import (
	"github.com/chris-joseph/golang-ecs/internal/api"
	"github.com/chris-joseph/golang-ecs/pkg/config"
	"github.com/chris-joseph/golang-ecs/pkg/data"
)

func main()  {
	cfg := config.New()
	db := data.NewMongoconnection(cfg)
	defer db.Disconnect()
	application := api.New(cfg,db.Client)
	application.Start()
}