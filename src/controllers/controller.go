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

type response map[string]interface{}

var validate = validator.New()

func getJSONMessage(Type, message string) map[string]interface{} {
	switch Type {
		case "e":
			return response{"error": message}
		default:
			return response{"message": message}
	}
}

func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Login hit!")
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "Bad Request"))
		return
	}

	if user.Email == "" || user.Password == "" {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "Both email and password are required"))
		return
	}

	existingUser := &models.User{}

	userProp := mgm.Coll(existingUser).FindOne(context.Background(), bson.M{"email": user.Email})

	if userProp.Err() != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "email or password is invalid"))
		return
	}

	var userFound models.User

	decodeErr := userProp.Decode(&userFound)
	if decodeErr != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "Internal server error"))
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {

		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "email or password is invalid"))
		return
	}

	res.WriteHeader(http.StatusOK)
	data := getJSONMessage("m", "User logged in")
	data["user"] = userFound
	json.NewEncoder(res).Encode(data)
}

// cr, _ := mgm.Coll(existingUser).Find(context.Background(), )
func Register(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Register hit!")
	res.Header().Set("Content-Type", "application/json")
	
	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "badRequest"))
		return
	}

	if validateError := validate.Struct(user); validateError != nil {
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(getJSONMessage("e", "All fields are required"))
		return
	}

	existingUser := &models.User{}
	if checkUserErr := mgm.Coll(existingUser).First(bson.M{"email": user.Email}, existingUser); checkUserErr != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if hashErr != nil {

			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(getJSONMessage("e", "internalServerError"))
			return
		}

		user.Password = string(hashedPassword)
		saveErr := mgm.Coll(&user).Create(&user)

		if saveErr != nil {
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(getJSONMessage("e", "internalServerError"))
			return 
		}

		res.WriteHeader(http.StatusOK)
		data := getJSONMessage("m", "user registered")
		data["user"] = user
		json.NewEncoder(res).Encode(data)
		return
	}

	res.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(res).Encode(getJSONMessage("e", "emailExists"))
}

func Home(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Home hit!")
}
