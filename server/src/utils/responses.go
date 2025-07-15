package utils

import (
	"encoding/json"
	"net/http"
)

type GlobalMap map[string]any

func GetJSONMessage(Type, message string) GlobalMap {
	switch Type {
		case "e":
			return GlobalMap{"error": message}
		default:
			return GlobalMap{"message": message}
	}
}

func Response (res http.ResponseWriter, responseMessage, status string, resData ...interface{}) {
	switch status {
	case "e":
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(GetJSONMessage("e", responseMessage))
	case "b":
		res.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(res).Encode(GetJSONMessage("e", responseMessage))
	case "u":
		res.WriteHeader(http.StatusUnauthorized)
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