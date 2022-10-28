package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"top-up-service/config"
	"top-up-service/helper"

	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
)

type TopUpRequest struct {
	ID           string  `json:"id"`
	Harga        int     `json:"harga"`
	HargaBuyback int     `json:"harga_buyback"`
	Norek        string  `json:"norek"`
	Gram         float64 `json:"gram"`
	RekeningId   string  `json:"rekening_id"`
}

func topup(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	tu := TopUpRequest{}
	decoder := json.NewDecoder(r.Body)
	reffID, _ := shortid.Generate()

	tu.ID = reffID
	if err := decoder.Decode(&tu); err != nil {
		helper.Respond(w, http.StatusBadRequest, true, reffID, err.Error())
		return
	}
	defer r.Body.Close()

	if tu.Gram > 0 {
		dot := fmt.Sprintf("%.4f", tu.Gram-math.Floor(tu.Gram))[5:]
		dotInt, err := strconv.Atoi(dot)
		if err != nil {
			helper.Respond(w, http.StatusUnprocessableEntity, true, reffID, err.Error())
		}
		if dotInt > 0 {
			helper.Respond(w, http.StatusUnprocessableEntity, true, reffID, "minimum top-up is multiply of 0.001")
			return
		}
	}
	err := c.Db.QueryRow("SELECT id FROM rekening where no_rek=$1", tu.Norek).Scan(&tu.RekeningId)
	if err != nil {
		log.Printf("Get Rekening Failed err : %v", err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, reffID, err.Error())
		return
	}
	err = c.Db.QueryRow("SELECT harga_topup, harga_buyback FROM harga where harga_topup=$1", tu.Harga).Scan(&tu.Harga, &tu.HargaBuyback)
	if err != nil {
		log.Printf("Get harga Failed err : %v", err.Error())
		helper.Respond(w, http.StatusUnprocessableEntity, true, reffID, err.Error())
		return
	}

	payloadBytes, err := json.Marshal(&tu)
	if err != nil {
		helper.Respond(w, http.StatusUnprocessableEntity, true, reffID, err.Error())
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
func gramValidation(gram float64) error {

	return nil
}
func main() {
	a := config.NewConfig()
	a.Rt.HandleFunc("/api/topup", topup).Methods("POST")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Serving top-up service at %v", config.URL())
	log.Fatal(srv.ListenAndServe())
}
