package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/chris-joseph/golang-ecs/pkg/config"
	"github.com/chris-joseph/golang-ecs/pkg/data"
	"github.com/chris-joseph/golang-ecs/pkg/domain"
	"github.com/chris-joseph/golang-ecs/pkg/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)


type IUserService interface{
	CreateAccount(user *domain.User) *models.Error
	Login(user *domain.User) (string,*models.Error)
}

type UserService struct{
	userProvider data.IUserProvider 
	cfg          *config.Settings
}

func NewUserService(cfg *config.Settings,userProvider data.IUserProvider)IUserService{
	return &UserService{
		userProvider: userProvider,
		cfg:cfg,
	}
}

func (u UserService) CreateAccount(user *domain.User) *models.Error{
		userExists,err:=u.userProvider.UsernameExists(user.UserName)

		if err != nil {
			return &models.Error{
				Code: 500,
				Name:"SERVER_ERROR",
				Message:"Someting went wrong",
				Error:err,
			}
		}
		if userExists{
			return &models.Error{
				Code: 400,
				Name:"USERNAME_TAKEN",
				Message:"Username already exists",
			}
		}
		user.ID=primitive.NewObjectID()
		hash,err:=hashPassword(user.Password)

		if err != nil {
			 return &models.Error{
				Code: 500,
				Name:"SERVER_ERROR",
				Message:"Someting went wrong",
				Error:err,
			}
		}

		user.Password = hash
		err = u.userProvider.CreateAccount(user)

		if err != nil {
			return &models.Error{
			   Code: 500,
			   Name:"SERVER_ERROR",
			   Message:"Someting went wrong",
			   Error:err,
		   }
	   }

		return nil
}


func (u UserService)Login(user *domain.User) (string,*models.Error){
	userFound,err:=u.userProvider.FindUserByName(user.UserName)

		if err != nil {
		fmt.Println(err)
			return "",&models.Error{
				Code: 500,
				Name:"SERVER_ERROR",
				Message:"Someting went wrong",
				Error:err,
			}
		}

		if userFound==nil{
			return "", &models.Error{
				Code: 400,
				Name:"INVALID_LOGIN",
				Message:"Username or Password is in valid",
			}
		}

		err = comparePassWordWithHash(user.Password,userFound.Password)

		if  err != nil{
			return "", &models.Error{
				Code: 400,
				Name:"INVALID_LOGIN",
				Message:"Username or Password is in valid",
			}
		}
		token,err:=u.CreateJWTToken(userFound.ID.Hex())

		if err != nil {
			fmt.Println(err)
			return "",&models.Error{
				Code: 500,
				Name:"SERVER_ERROR",
				Message:"Someting went wrong",
				Error:err,
			}
		}

		return token,nil

}

func hashPassword(password string) (string,error) {
	passwordBytes:= []byte(password)

	hashedPassword,err:=bcrypt.GenerateFromPassword(passwordBytes,12)

	if err != nil {
		return "",errors.Wrap(err,"Error creating passsword")
	}

	return string(hashedPassword),nil
}

func comparePassWordWithHash(password string,hash string) error {
	
	passwordBytes := []byte(password)
	hashBytes := []byte(hash)

	err:=bcrypt.CompareHashAndPassword(hashBytes,passwordBytes)

	return errors.Wrap(err,"Error comparing passsword and hash")

}

func (u UserService)CreateJWTToken(userID string) (string,error ) {

	token := jwt.New(jwt.SigningMethodHS256)

	expiresIn,err:=strconv.ParseInt(u.cfg.JwtExpires,10,64)

	if err != nil {
		return "",errors.Wrap(err,"Error parsing int")
	}
	expiration:=time.Duration(int64(time.Minute)* expiresIn)

	claims:= token.Claims.(jwt.MapClaims)

	claims["id"]=userID

	claims["exp"]=time.Now().Add(expiration).Unix()

	t,err:=token.SignedString([]byte(u.cfg.JwtSecret))


	if err != nil {
		return "", errors.Wrap(err,"Error signing jwt token")
	}
	return t,nil
}