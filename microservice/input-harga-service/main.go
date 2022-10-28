package main

import (
	"encoding/json"
	"fmt"
	"input-harga-service/config"
	"input-harga-service/helper"
	"log"
	"net/http"

	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/kafka-go"
	"github.com/teris-io/shortid"
)

type HargaRequest struct {
	ID           string `json:"id"`
	HargaTopup   int    `json:"harga_topup"`
	HargaBuyback int    `json:"harga_buyback"`
	AdminID      string `json:"admin_id"`
}

func createHarga(w http.ResponseWriter, r *http.Request) {
	c := config.NewConfig()
	params := HargaRequest{}
	decoder := json.NewDecoder(r.Body)
	reffID, _ := shortid.Generate()
	params.ID = reffID
	if err := decoder.Decode(&params); err != nil {
		helper.Respond(w, http.StatusBadRequest, true, reffID, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	payloadBytes, err := json.Marshal(&params)
	if err != nil {
		helper.Respond(w, http.StatusBadRequest, true, reffID, err.Error())
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
	a.Rt.HandleFunc("/api/input-harga", createHarga).Methods("POST")
	srv := &http.Server{
		Handler:      a.Rt,
		Addr:         config.URL(),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println(config.URL())
	log.Fatal(srv.ListenAndServe())
}
