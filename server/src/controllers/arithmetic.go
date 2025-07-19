package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"

	"github.com/cakezero/go-server/src/models"
	"github.com/cakezero/go-server/src/utils"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

var ctx = context.Background()

func ArithmeticHistory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	id := req.URL.Query().Get("id")
	fmt.Println("User ID history:", id)

	var history []models.Arithmetic

	historyFetch, findErr := mgm.Coll(&models.Arithmetic{}).Find(ctx, bson.M{"user": id})
	if findErr != nil {
		fmt.Println(findErr.Error())
		utils.Response(res, "No arithmetic has been performed", "")
		return
	}

	if fetchErr := historyFetch.All(ctx, &history); fetchErr != nil {
		utils.Response(res, "Error fetching history", "e")
		return
	}

	utils.Response(res, "arithmetic history fetched!", "", history)
}

func PerformAction(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	var userArithmetic models.Arithmetic
	type computation struct {
		First string
		Second string
		Action string
		Id string
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

	var equation string

	if check := answer == float64(int64(answer)); check {
    equation = fmt.Sprintf("%d %s %d 🟰 %d", int(firstNumber), symbol, int(secondNumber), int(answer))
	} else {
		equation = fmt.Sprintf("%.1f %s %.1f 🟰 %.03f", firstNumber, symbol, secondNumber, answer)
	}

	id := userAction.Id

	if id != "" {
		userArithmetic.Equation = equation
		userArithmetic.Action = action
		userArithmetic.User = userAction.Id

		saveErr := mgm.Coll(&userArithmetic).Create(&userArithmetic)

		if saveErr != nil {
			utils.Response(res, "Error saving arithmetic calculation", "e")
			return
		}

		utils.Response(res, "arithmetic action performed", "", answer)
		return
	}

	data := utils.GlobalMap{
		"answer": answer,
		"equation": equation,
		"action": action,
	}	

	utils.Response(res, "arithmetic action performed", "", data)
}
