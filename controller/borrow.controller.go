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

var BorrowBuku = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}

	borrow := &models.RequestPinjam{}
	err = json.NewDecoder(r.Body).Decode(borrow)

	resp, err := borrow.PinjamBuku(GetDb(), uint(id), w)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
	conn, err := GetDb().DB()
	if err != nil {
		defer conn.Close()
	}
}

var ReturnBook = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])

	if err != nil {
		fmt.Println(err)
	}

	returnBook := &models.ReturnBook{}
	err = json.NewDecoder(r.Body).Decode(returnBook)
	resp, err := returnBook.ReturnBook(GetDb(), uint(id), w)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Input")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
	conn, err := GetDb().DB()
	if err != nil {
		defer conn.Close()
	}

}
var ListBorrow = func(w http.ResponseWriter, r *http.Request) {

	listBorrow := &models.OrderDetail{}
	resp, err := listBorrow.ListBorrow(GetDb(), w)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Not Found")
	}
	//
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
	conn, err := GetDb().DB()
	if err != nil {
		defer conn.Close()
	}
}

var ListReturnBook = func(w http.ResponseWriter, r *http.Request) {

	listReturn := &models.OrderDetail{}
	resp, err := listReturn.ListReturnBook(GetDb(), w)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Not Found")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
	conn, err := GetDb().DB()
	if err != nil {
		defer conn.Close()
	}
}
