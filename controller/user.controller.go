package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"onboarding/helpers"
	"onboarding/models"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var dbConfig *gorm.DB

func InitiateDb(db *gorm.DB) {
	dbConfig = db
}

func GetDb() *gorm.DB {
	return dbConfig
}

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	newUser := &models.User{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	_, resp := newUser.Create(GetDb())
	helpers.ResponseWithJson(w, 200, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	userLogin := &models.User{}
	err := json.NewDecoder(r.Body).Decode(userLogin)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	resp, err := userLogin.Login(GetDb(), userLogin.Mobile, userLogin.Email, userLogin.Password)
	helpers.ResponseWithJson(w, 200, resp)
}

var GetUserById = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}

	user := &models.User{}

	resp, err := user.GetUserById(GetDb(), uint(id))
	helpers.ResponseWithJson(w, 200, resp)

}

var GetUsers = func(w http.ResponseWriter, r *http.Request) {
	users := &models.User{}

	resp, err := users.GetUsers(GetDb())
	if err != nil {
		fmt.Println(err)
	}
	helpers.ResponseWithJson(w, 200, resp)
}

var UpdateUsers = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Invalid Request")
	}
	resp, _ := user.UpdateUsers(GetDb(), uint(id))
	if resp != nil {
		helpers.ResponseWithJson(w, 200, resp)
	}
}

var UpdatePhoto = func(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	resp, err := user.UpdatePhoto(GetDb(), uint(id), r)
	if resp != nil {
		helpers.ResponseWithJson(w, 200, resp)
	}
}

var DeleteUser = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		fmt.Println(err)
	}
	user := &models.User{}
	resp, err := user.DeleteUserByID(GetDb(), uint(id))
	if resp != nil {
		helpers.ResponseWithJson(w, 200, resp)
	}

}

// func getDB() *gorm.DB {
// 	conn := models.GetDB()
// 	dbSQL, ok := conn.DB()
// 	if ok != nil {
// 		defer dbSQL.Close()
// 	}
// 	return conn
// }
