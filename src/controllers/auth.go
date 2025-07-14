package controllers

import (
	"context"
	"fmt"
	"net/http"

	"encoding/json"

	"github.com/cakezero/go-server/src/models"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
	"golang.org/x/crypto/bcrypt"
)

type GlobalMap map[string]any

var validate = validator.New()

func GetJSONMessage(Type, message string) GlobalMap {
	switch Type {
		case "e":
			return GlobalMap{"error": message}
		default:
			return GlobalMap{"message": message}
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Login hit!")
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		Response(res, "Data received from request is not appropriate", "b")
		return
	}

	if user.Email == "" || user.Password == "" {
		Response(res, "Both email and password are required", "b")
		return
	}

	userProp := mgm.Coll(&models.User{}).FindOne(context.Background(), bson.M{"email": user.Email})

	if userProp.Err() != nil {
		Response(res, "email or password is invalid", "b")
		return
	}

	var userFound models.User

	decodeErr := userProp.Decode(&userFound)
	if decodeErr != nil {
		Response(res, "internal server error", "e")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		Response(res, "email or password is invalid", "b")
		return
	}

	Response(res, "user logged in", "", userFound)
}

func Register(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Register hit!")
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		Response(res, "Request data received was not appropriate", "b")
		return
	}

	if validateError := validate.Struct(user); validateError != nil {
		Response(res, "All fields are required", "b")
		return
	}

	if length := len(user.Password); length < 8 {
		Response(res, "Password must be gte 8", "b")
		return
	}

	existingUser := &models.User{}
	if checkUserErr := mgm.Coll(existingUser).First(bson.M{"email": user.Email}, existingUser); checkUserErr != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if hashErr != nil {
			Response(res, "Internal server error", "e")
			return
		}

		user.Password = string(hashedPassword)
		saveErr := mgm.Coll(&user).Create(&user)

		if saveErr != nil {
			Response(res, "User not saved", "e")
			return
		}

		Response(res, "User saved", "", user)
		return
	}

	Response(res, "Email exists", "b")
}

func Home(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Home hit!")
}
