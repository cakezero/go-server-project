package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/cakezero/go-server/src/middlewares"
	"github.com/cakezero/go-server/src/models"
	"github.com/cakezero/go-server/src/utils"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

var Ctx = context.Background()

func ArithmeticHistory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id := req.Context().Value(middlewares.IdKey).(string)
	fmt.Printf("id: %s\n", id)
	var history []models.Arithmetic

	historyFetch, findErr := mgm.Coll(&models.Arithmetic{}).Find(Ctx, bson.M{"user": id})
	if findErr != nil {
		fmt.Println(findErr.Error())
		utils.Response(res, "No arithmetic has been performed", "")
		return
	}

	if fetchErr := historyFetch.All(Ctx, &history); fetchErr != nil {
		utils.Response(res, "Error fetching history", "e")
		return
	}

	utils.Response(res, "arithmetic history fetched!", "", history)
}

func PerformAction(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := req.Context().Value(middlewares.IdKey).(string)
	fmt.Printf("id: %s\n", id)
	var userArithmetic models.Arithmetic
	type computation struct {
		First string
		Second string
		Action string
	}

	var userAction computation

	if decodeErr := json.NewDecoder(req.Body).Decode(&userAction); decodeErr != nil {
		utils.Response(res, "Data received from request is not appropriate", "e")
		return
	}

	firstNumber, firstErr := strconv.ParseFloat(userAction.First, 64)
	if firstErr != nil {
		utils.Response(res, "Only integer strings are allowed", "b")
		return
	}

	secondNumber, secondErr := strconv.ParseFloat(userAction.Second, 64)
	if secondErr != nil {
		utils.Response(res, "Only integer strings are allowed", "b")
		return
	}

	var answer float64
	var symbol string
	action := userAction.Action

	switch action {
	case "add":
		answer = firstNumber + secondNumber
		symbol = "➕"
	case "subtract":
		answer = firstNumber - secondNumber
		symbol = "➖"
	case "multiply":
		answer = firstNumber * secondNumber
		symbol = "✖️"
	case "divide":
		answer = math.Round((firstNumber / secondNumber) * 1000) /1000
		symbol = "➗"
	default: 
		utils.Response(res, "Specify arithmetic action", "b")
		return
	}

	userArithmetic.Equation = fmt.Sprintf("%.0f %s %.0f 🟰 %.0f", firstNumber, symbol, secondNumber, answer)
	userArithmetic.Action = action
	userArithmetic.User = id

	saveErr := mgm.Coll(&userArithmetic).Create(&userArithmetic)

	if saveErr != nil {
		utils.Response(res, "Error saving arithmetic calculation", "e")
		return
	}

	utils.Response(res, "arithmetic action performed", "", answer)
}
