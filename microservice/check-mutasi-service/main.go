package main

import (
	"check-saldo-service/config"
	"check-saldo-service/helper"
	"encoding/json"
	"log"
	"net/http"

	"time"

	_ "github.com/lib/pq"
)

type MutasiRequest struct {
	NoRek     string `json:"no_rek"`
	StartDate int    `json:"start_date"`
	EndDate   int    `json:"end_date"`
}

type MutasiResponse struct {
	Date         int     `json:"date"`
	Type         string  `json:"type"`
	Gram         float64 `json:"gram"`
	HargaTopup   int     `json:"harga_topup"`
	HargaBuyback int     `json:"harga_buyback"`
	Saldo        float64 `json:"saldo"`
}

func checkMutasi(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	var (
		mReq       MutasiRequest
		mutasi     MutasiResponse
		mResp      []MutasiResponse
		rekeningId string
	)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&mReq); err != nil {
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}
	defer r.Body.Close()

	err := c.Db.QueryRow("SELECT id FROM rekening where no_rek=$1", mReq.NoRek).
		Scan(&rekeningId)
	if err != nil {
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}

	rows, err := c.Db.Query("SELECT created_at, type, gram, harga_topup, harga_buyback, saldo FROM transaksi where created_at>$1 and created_at<$2 and rekening_id=$3", mReq.StartDate, mReq.EndDate, rekeningId)

	if err != nil {
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&mutasi.Date, &mutasi.Type, &mutasi.Gram, &mutasi.HargaTopup, &mutasi.HargaBuyback, &mutasi.Saldo); err != nil {
			helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
			return
		}
		mResp = append(mResp, mutasi)
	}
	helper.Respond(w, http.StatusOK, false, "", mResp)
}
func main() {
	a := config.NewConfig()
	a.Rt.HandleFunc("/api/mutasi", checkMutasi).Methods("GET")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(config.URL())
	log.Fatal(srv.ListenAndServe())
}
