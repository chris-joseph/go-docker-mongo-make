package main

import ("github.com/chris-joseph/golang-ecs/internal/api")

func main()  {
	application:=api.New()
	application.Start()
}