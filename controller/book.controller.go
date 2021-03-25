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

var CreateBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}
	newBook := &models.Book{}
	err = json.NewDecoder(r.Body).Decode(newBook)

	resp, _ := newBook.Create(conn, uint(id), w)
	helpers.ResponseWithJson(w, 200, resp)
}

var GetBookByID = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}
	book := &models.Book{}
	resp, err := book.GetBookById(conn, uint(id))
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)

}

var GetAllBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	books := &models.Book{}

	resp, err := books.GetAllBook(conn)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}

var UpdateBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	params := mux.Vars(r)
	idUser, err := strconv.Atoi(params["id"])
	idBook, err := strconv.Atoi(r.URL.Query().Get("idBook"))

	if err != nil {
		fmt.Println(err)
	}
	book := &models.Book{}
	err = json.NewDecoder(r.Body).Decode(book)
	resp, err := book.UpdateBook(conn, uint(idBook), uint(idUser))
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}
var NewestBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	books := &models.Book{}

	resp, err := books.NewestBook(conn)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}

var PopularBook = func(w http.ResponseWriter, r *http.Request) {
	conn := getDB()
	histories := &models.History{}

	resp, err := histories.PopulerBook(conn)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	helpers.ResponseWithJson(w, http.StatusAccepted, resp)
}
