package main

import (
	"cryptoprice/pkg"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	pkg.InitDB()
	defer pkg.CloseDB()

	pkg.InitRabbitMQ()
	defer pkg.CloseRabbitMQ()

	router := mux.NewRouter()

	router.HandleFunc("/alerts/create", pkg.CreatePriceAlert).Methods("POST")
	router.HandleFunc("/alerts/delete/{id}", pkg.DeletePriceAlert).Methods("DELETE")
	router.HandleFunc("/alerts", pkg.GetPriceAlerts).Methods("GET").Queries("page", "{page:[0-9]+}", "status", "{status}")

	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			pkg.Authenticate(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				next.ServeHTTP(w, r)
			})).ServeHTTP(w, r)
		})
	})

	http.Handle("/", router)

	addr := ":8080"
	fmt.Printf("Server is running on %s...\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		panic(err)
	}
}
