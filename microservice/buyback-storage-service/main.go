package main

import (
	"buyback-storage-service/config"
	"context"
	"encoding/json"
	"log"
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

func main() {
	c := config.NewConfig()
	ctx := context.Background()
	for {
		m, err := c.Kr.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))
		var b Buyback
		if err := json.Unmarshal(m.Value, &b); err != nil {
			log.Printf("unmarshall data error : %s", err.Error())
		}

		err = c.Db.QueryRow("UPDATE rekening SET saldo=$1 where id=$2", b.Saldo, b.RekID).Err()
		if err != nil {
			log.Printf("Get harga Failed err : %v", err.Error())
			return
		}

		err = c.Db.QueryRow("INSERT INTO transaksi(id, type, rekening_id, saldo, gram, harga_topup, harga_buyback) VALUES($1, $2, $3, $4, $5, $6, $7)",
			b.ID, "buyback", b.RekID, b.Saldo, b.Gram, b.HargaTopup, b.Harga).Err()
		if err != nil {
			log.Printf("Get harga Failed err : %v", err.Error())
			return
		}

		if err := c.Kr.CommitMessages(ctx, m); err != nil {
			log.Printf("CommitMessage Failed error : %s", err.Error())
		}
	}
}
