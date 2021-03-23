package main

import (
	"log"
	"net/http"
	"onboarding/controller"
	"onboarding/middleware"
	"onboarding/models"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	models.Init()
	router := mux.NewRouter()

	router.HandleFunc("/", index)

	subRouter := router.PathPrefix("/user").Subrouter()
	subProtectedRouter := router.PathPrefix("/user").Subrouter()

	subRouter.HandleFunc("/v1/signup", controller.CreateAccount).Methods("POST")
	subRouter.HandleFunc("/v1/login", controller.Authenticate).Methods("POST")
	//router.HandleFunc("/api/v1/user/signup", controller.CreateAccount).Methods("POST")
	//router.HandleFunc("/api/v1/user/login", controller.Authenticate).Methods("POST")
	subRouter.HandleFunc("/users", controller.GetUsers).Methods("GET")
	subProtectedRouter.Use(middleware.JwtVerifyToken)
	subProtectedRouter.HandleFunc("/v1/user/{id}", controller.GetUserById).Methods("GET")
	subProtectedRouter.HandleFunc("/v1/user/{id}", controller.UpdateUsers).Methods("PUT")
	subProtectedRouter.HandleFunc("/v1/user/photo/{id}", controller.UpdatePhoto).Methods("PUT")
	subProtectedRouter.HandleFunc("/v1/user/{id}", controller.DeleteUser).Methods("DELETE")
	subProtectedRouter.HandleFunc("/v1/book/{id}", controller.CreateBook).Methods("POST")
	// subProtectedRouter.HandleFunc("/v1/book/{id}", controller.CreateBook).Methods("POST")
	subProtectedRouter.HandleFunc("/v1/book/{id}", controller.GetBookByID).Methods("GET")

	subProtectedRouter.HandleFunc("/v1/books", controller.GetAllBook).Methods("GET")
	subProtectedRouter.HandleFunc("/v1/book/{id}", controller.UpdateBook).Methods("PUT")
	subProtectedRouter.HandleFunc("/v1/booknewest", controller.NewestBook).Methods("GET")
	subProtectedRouter.HandleFunc("/v1/stock/{id}", controller.CreateStock).Methods("POST")
	subProtectedRouter.HandleFunc("/v1/popularBook", controller.PopularBook).Methods("GET")

	subProtectedRouter.HandleFunc("/v1/borrow/{id}", controller.BorrowBuku).Methods("POST")
	subProtectedRouter.HandleFunc("/v1/return/{id}", controller.ReturnBook).Methods("POST")

	subProtectedRouter.HandleFunc("/v1/listBorrow", controller.ListBorrow).Methods("GET")

	router.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Server is running"))
}
