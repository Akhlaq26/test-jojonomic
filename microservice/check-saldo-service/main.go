package main

import (
	"check-saldo-service/config"
	"check-saldo-service/helper"
	"log"
	"net/http"

	"time"

	_ "github.com/lib/pq"
)

type Saldo struct {
	NoRek int     `json:"no_rek"`
	Saldo float64 `json:"saldo"`
}

func checkHarga(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	var s Saldo
	err := c.Db.QueryRow("SELECT no_rek, saldo FROM rekening").Scan(&s.NoRek, &s.Saldo)
	if err != nil {
		helper.Respond(w, http.StatusBadRequest, true, "", err.Error())
		return
	}
	helper.Respond(w, http.StatusOK, false, "", s)
}
func main() {
	a := config.NewConfig()
	a.Rt.HandleFunc("/api/saldo", checkHarga).Methods("GET")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(config.URL())
	log.Fatal(srv.ListenAndServe())
}
