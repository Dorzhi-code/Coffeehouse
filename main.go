package main

import (
	databse "coffe/database"
	"coffe/handler"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type homeHandler struct{}

func (h *homeHandler) ServerHTTP(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("This is my home page"))
}

func main(){
	db, err := databse.Connect("./database/coffee.db")
	
	if(err != nil){
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	customerHandler := handler.CustomerHandler{DB:db}

	home := homeHandler{}
	
	router.HandleFunc("/", home.ServerHTTP)
	router.HandleFunc("/customers", customerHandler.Create).Methods("POST")
	router.HandleFunc("/customers", customerHandler.RetrieveAll).Methods("GET")
	router.HandleFunc("/customers/{id}", customerHandler.Retrieve).Methods("GET")
	router.HandleFunc("/customers/{id}", customerHandler.Update).Methods("PUT")
	router.HandleFunc("/customers/{id}", customerHandler.Delete).Methods("DELETE")

	http.ListenAndServe(":8010", router)
}