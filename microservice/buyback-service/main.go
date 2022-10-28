package main

import (
	"buyback-service/config"
	"buyback-service/helper"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
)

type Buyback struct {
	ID         string  `json:"id"`
	Gram       float64 `json:"gram"`
	Harga      int     `json:"harga"`
	NoRek      string  `json:"no_rek"`
	HargaTopup int     `json:"harga_topup"`
	RekID      string  `json:"rek_id"`
	Saldo      float64 `json:"saldo"`
}

func topup(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	b := Buyback{}
	decoder := json.NewDecoder(r.Body)
	reffID, _ := shortid.Generate()

	b.ID = reffID
	if err := decoder.Decode(&b); err != nil {
		helper.Respond(w, http.StatusBadRequest, true, "", err.Error())
		return
	}
	defer r.Body.Close()
	saldo := float64(0)
	err := c.Db.QueryRow("SELECT id, saldo FROM rekening where no_rek=$1", b.NoRek).Scan(&b.RekID, &saldo)
	if err != nil {
		log.Printf("Get Saldo Failed err : %v", err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}
	if b.Gram < saldo {
		err = errors.New("Gram must be bigger than saldo")
		log.Print(err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}
	b.Saldo = saldo - b.Gram
	err = c.Db.QueryRow("SELECT harga_topup FROM harga where harga_buyback=$1", b.Harga).Scan(&b.HargaTopup)
	if err == sql.ErrNoRows {
		log.Print(err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", "Harga not Found")
		return
	}
	if err != nil {
		log.Printf("Get harga Failed err : %v", err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}

	payloadBytes, err := json.Marshal(&b)
	if err != nil {
		helper.Respond(w, http.StatusUnprocessableEntity, true, "", err.Error())
		return
	}

	c.Kafka.SetWriteDeadline(time.Now().Add(10 * time.Second))
	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("address-%s", r.RemoteAddr)),
		Value: payloadBytes,
	}
	_, err = c.Kafka.WriteMessages(msg)
	if err != nil {
		log.Println(err.Error())
		helper.Respond(w, http.StatusBadRequest, true, reffID, "Kafka not ready")
		return
	}

	helper.Respond(w, http.StatusOK, false, reffID, nil)
}
func main() {
	a := config.NewConfig()
	a.Rt.HandleFunc("/api/buyback", topup).Methods("POST")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Serving top-up service at %v", config.URL())
	log.Fatal(srv.ListenAndServe())
}
