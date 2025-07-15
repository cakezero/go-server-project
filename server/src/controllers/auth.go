package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/cakezero/go-server/src/models"
	"github.com/cakezero/go-server/src/utils"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()


func Login(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Login hit!")
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.Response(res, "Data received from request is not appropriate", "b")
		return
	}

	if user.Email == "" || user.Password == "" {
		utils.Response(res, "Both email and password are required", "b")
		return
	}

	userProp := mgm.Coll(&models.User{}).FindOne(context.Background(), bson.M{"email": user.Email})

	if userProp.Err() != nil {
		utils.Response(res, "email or password is invalid", "b")
		return
	}

	var userFound models.User

	decodeErr := userProp.Decode(&userFound)
	if decodeErr != nil {
		utils.Response(res, "internal server error", "e")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password)); err != nil {
		utils.Response(res, "email or password is invalid", "b")
		return
	}

	accessToken, refreshToken, tokenFetchErr := utils.GenerateJWTs(user.ID.String())

	if tokenFetchErr != nil {
		utils.Response(res, "Error creating access/refresh token", "e")
		return
	}

	http.SetCookie(res, &http.Cookie{
		Name: "refresh_token",
		Value: refreshToken,
		HttpOnly: true,
		Secure: true,
		Path: "/",
		MaxAge: 7 * 24 * 60 * 60,
	})

	data := utils.GlobalMap{
		"user": userFound,
		"token": accessToken,
	}

	utils.Response(res, "user logged in", "", data)
}

func Register(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Register hit!")
	res.Header().Set("Content-Type", "application/json")

	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		utils.Response(res, "Request data received was not appropriate", "b")
		return
	}

	if validateError := validate.Struct(user); validateError != nil {
		utils.Response(res, "All fields are required", "b")
		return
	}

	if length := len(user.Password); length < 8 {
		utils.Response(res, "Password must be gte 8", "b")
		return
	}

	existingUser := &models.User{}
	if checkUserErr := mgm.Coll(existingUser).First(bson.M{"email": user.Email}, existingUser); checkUserErr != nil {
		hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if hashErr != nil {
			utils.Response(res, "Internal server error", "e")
			return
		}

		user.Password = string(hashedPassword)
		saveErr := mgm.Coll(&user).Create(&user)

		if saveErr != nil {
			utils.Response(res, "User not saved", "e")
			return
		}

		accessToken, refreshToken, tokenFetchErr := utils.GenerateJWTs(user.ID.String())

		if tokenFetchErr != nil {
			utils.Response(res, "Error creating access/refresh token", "e")
			return
		}

		http.SetCookie(res, &http.Cookie{
			Name: "refresh_token",
			Value: refreshToken,
			HttpOnly: true,
			Secure: true,
			Path: "/",
			MaxAge: 7 * 24 * 60 * 60,
		})

		data := utils.GlobalMap{
			"user": user,
			"token": accessToken,
		}

		utils.Response(res, "User saved", "", data)
		return
	}

	utils.Response(res, "Email exists", "b")
}

func RefreshTokenHandler(res http.ResponseWriter, req *http.Request) {
	refresh_cookie, err := req.Cookie("refresh_token")

	if err != nil {
		utils.Response(res, "Refresh token not found", "u")
		return
	}

	userId, _ := utils.DecodeJWT(res, refresh_cookie.Value)

	if userId == "" {
		return
	}

	accessToken, _, err := utils.GenerateJWTs(userId)
	if err != nil {
		utils.Response(res, "Failed to generate JWT", "e")
		return
	}

	data := utils.GlobalMap{
		"token": accessToken,
	}

	utils.Response(res, "Access token refreshed", "", data)
}

func Logout(res http.ResponseWriter, req *http.Request) {
	fmt.Println("Logout hit!")

	authHeader := req.Header.Get("Authorizaton")
	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		utils.Response(res, "Invalid auth header format", "u")
		return
	}

	accessToken := parts[1]
	_, claims := utils.DecodeJWT(res, accessToken)

	exp := int64(0)
	if expFloat, ok := claims["exp"].(float64); ok {
		exp = int64(expFloat)
	}

	if exp > 0 {
		ttl := time.Until(time.Unix(exp, 0))
		err := utils.GetRedisClient().Set(ctx, accessToken, "revoked", ttl).Err()

		if err != nil {
			utils.Response(res, "Error revoking token", "e")
			return
		}
	}

	http.SetCookie(res, &http.Cookie{
		Name: "refresh_token",
		Value: "",
		MaxAge: -1,
		Path: "/",
		HttpOnly: true,
		Secure: true,
		Expires: time.Unix(0, 0),
	})

	utils.Response(res, "User logged out", "")
}
