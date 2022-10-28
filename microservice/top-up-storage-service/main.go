package main

import (
	"context"
	"encoding/json"
	"log"
	"top-up-storage-service/config"
)

type TopUpRequest struct {
	ID           string  `json:"id"`
	Harga        int     `json:"harga"`
	HargaBuyback int     `json:"harga_buyback"`
	Norek        string  `json:"norek"`
	Gram         float64 `json:"gram"`
	RekeningId   string  `json:"rekening_id"`
}

func main() {
	c := config.NewConfig()
	ctx := context.Background()
	for {
		m, err := c.Kr.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))
		var tu TopUpRequest
		if err := json.Unmarshal(m.Value, &tu); err != nil {
			log.Printf("unmarshall data error : %s", err.Error())
		}
		saldo := float64(0)
		err = c.Db.QueryRow(
			"SELECT saldo FROM rekening where id=$1",
			tu.RekeningId).Scan(&saldo)

		if err != nil {
			log.Printf("Get Saldo Failed error : %s", err.Error())
		}
		saldo = saldo + tu.Gram

		err = c.Db.QueryRow(
			"INSERT INTO top_up(id, harga, gram, rekening_id) VALUES($1, $2, $3, $4)",
			tu.ID, tu.Harga, tu.Gram, tu.RekeningId).Err()

		if err != nil {
			log.Printf("Create Top up Failed error : %s", err.Error())
		}

		err = c.Db.QueryRow(
			"UPDATE rekening SET saldo=$1 where id=$2",
			saldo, tu.RekeningId).Err()
		if err != nil {
			log.Printf("Update Rekening Failed error : %s", err.Error())
		}

		err = c.Db.QueryRow(
			"INSERT INTO transaksi(id, type, top_up_id, rekening_id, saldo, gram, harga_topup, harga_buyback) VALUES($1, $2, $3, $4, $5, $6, $7, &8)",
			tu.ID, "top_up", tu.ID, tu.RekeningId, saldo, tu.Gram, tu.Harga).Err()

		if err != nil {
			log.Printf("Create Transaction Failed error : %s", err.Error())
		}

		if err := c.Kr.CommitMessages(ctx, m); err != nil {
			log.Printf("CommitMessage Failed error : %s", err.Error())
		}

	}
}
