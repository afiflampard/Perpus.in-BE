package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"onboarding/helpers"
	"onboarding/models"
	"strconv"

	"github.com/gorilla/mux"
)

var CreateStock = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idUser, err := strconv.Atoi(params["id"])
	idBook, err := strconv.Atoi(r.URL.Query().Get("idBook"))
	if err != nil {
		fmt.Println(err)
	}
	stock := &models.Stock{}
	err = json.NewDecoder(r.Body).Decode(stock)

	resp, err := stock.Create(GetDb(), w, uint(idUser), uint(idBook))

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
	conn, err := GetDb().DB()
	if err != nil {
		defer conn.Close()
	}

}
