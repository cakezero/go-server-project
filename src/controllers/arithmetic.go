package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"github.com/cakezero/go-server/src/models"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"math"
)

var Ctx = context.Background()

func Response (res http.ResponseWriter, responseMessage, status string, resData ...interface{}) {
	switch status {
	case "e":
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(GetJSONMessage("e", responseMessage))
	case "b":
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(GetJSONMessage("e", responseMessage))
	default:
		if len(resData) > 0 {
			data := GetJSONMessage("", responseMessage)
			data["data"] = resData[0]
			res.WriteHeader(http.StatusOK)
			json.NewEncoder(res).Encode(data)
		} else {
			res.WriteHeader(http.StatusOK)
			json.NewEncoder(res).Encode(GetJSONMessage("", responseMessage))
		}
	}
}

func ArithmeticHistory(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	id := req.URL.Query().Get("id")
	var history []models.Arithmetic

	historyFetch, findErr := mgm.Coll(&models.Arithmetic{}).Find(Ctx, bson.M{"user": id})
	if findErr != nil {
		fmt.Println(findErr.Error())
		Response(res, "No arithmetic has been performed", "")
		return
	}

	if fetchErr := historyFetch.All(Ctx, &history); fetchErr != nil {
		Response(res, "Error fetching history", "e")
		return
	}

	Response(res, "arithmetic history fetched!", "", history)
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
		Response(res, "Data received from request is not appropriate", "e")
		return
	}

	firstNumber, firstErr := strconv.ParseFloat(userAction.First, 64)
	if firstErr != nil {
		Response(res, "Only integer strings are allowed", "b")
		return
	}

	secondNumber, secondErr := strconv.ParseFloat(userAction.Second, 64)
	if secondErr != nil {
		Response(res, "Only integer strings are allowed", "b")
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
		Response(res, "Specify arithmetic action", "b")
		return
	}

	userArithmetic.Equation = fmt.Sprintf("%.0f %s %.0f 🟰 %.0f", firstNumber, symbol, secondNumber, answer)
	userArithmetic.Action = action
	userArithmetic.User = userAction.Id

	saveErr := mgm.Coll(&userArithmetic).Create(&userArithmetic)

	if saveErr != nil {
		Response(res, "Error saving arithmetic calculation", "e")
		return
	}

	Response(res, "arithmetic action performed", "", answer)
}
