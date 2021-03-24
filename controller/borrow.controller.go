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
	conn := getDB()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	idBook := r.URL.Query().Get("idBook")
	if err != nil {
		fmt.Println(err)
	}

	borrow := &models.RequestPinjam{}
	err = json.NewDecoder(r.Body).Decode(borrow)

	resp, err := borrow.PinjamBuku(conn, uint(id), idBook, w)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}

var ReturnBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	idBook := r.URL.Query().Get("idBook")
	if err != nil {
		fmt.Println(err)
	}

	returnBook := &models.ReturnBook{}
	err = json.NewDecoder(r.Body).Decode(returnBook)
	resp, err := returnBook.ReturnBook(conn, idBook, uint(id), w)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Input")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)

}
var ListBorrow = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()

	listBorrow := &models.OrderDetail{}
	resp, err := listBorrow.ListBorrow(conn, w)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Not Found")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}

var ListReturnBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()

	listReturn := &models.OrderDetail{}
	resp, err := listReturn.ReturnBook(conn, w)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Not Found")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}
