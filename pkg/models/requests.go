package models

import (
	"github.com/chris-joseph/golang-ecs/pkg/domain"
	"github.com/labstack/echo/v4"
)

type RegisterRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}


type LoginRequest struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

func ValidateRegisterRequest(c echo.Context) (*domain.User,*Error) {
	registerRequest := new(RegisterRequest)
	if err := c.Bind(registerRequest);err!=nil {
		return nil,BindError()
	}
	var validationErrors [] string

	if len(registerRequest.Password) < 8 {
		validationErrors = append(validationErrors, "Password must be 8 characters")
		
	}

	if len(registerRequest.UserName) < 3 {
		validationErrors = append(validationErrors, "Username must be longer than 2 characters")
		
	}

	if len(validationErrors) > 0 {
		return nil,ValidationError(validationErrors)
		
	}
	return &domain.User{
		UserName: registerRequest.UserName,
		Password: registerRequest.Password,
	},nil


}


func ValidateLoginRequest(c echo.Context) (*domain.User,*Error) {
	loginRequest := new(LoginRequest)
	if err := c.Bind(loginRequest);err!=nil {
		return nil,BindError()
	}
	var validationErrors [] string

	if len(loginRequest.Password) < 8 {
		validationErrors = append(validationErrors, "Password must be 8 characters")
		
	}

	if len(loginRequest.UserName) < 3 {
		validationErrors = append(validationErrors, "Username must be longer than 2 characters")
		
	}

	if len(validationErrors) > 0 {
		return nil,ValidationError(validationErrors)
		
	}
	return &domain.User{
		UserName: loginRequest.UserName,
		Password: loginRequest.Password,
	},nil


}
