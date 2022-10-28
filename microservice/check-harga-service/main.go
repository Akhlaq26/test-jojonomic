package main

import (
	"input-harga-service/config"
	"input-harga-service/helper"
	"log"
	"net/http"

	"time"

	_ "github.com/lib/pq"
)

type HargaResponse struct {
	HargaTopup   int `json:"harga_topup"`
	HargaBuyback int `json:"harga_buyback"`
}

func checkHarga(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	var h HargaResponse
	err := c.Db.QueryRow("SELECT harga_topup, harga_buyback FROM harga").Scan(&h.HargaTopup, &h.HargaBuyback)
	if err != nil {
		helper.Respond(w, http.StatusBadRequest, true, "", err.Error())
		return
	}
	helper.Respond(w, http.StatusOK, false, "", h)
}
func main() {
	a := config.NewConfig()
	a.Rt.HandleFunc("/api/check-harga", checkHarga).Methods("GET")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(config.URL())
	log.Fatal(srv.ListenAndServe())
}
