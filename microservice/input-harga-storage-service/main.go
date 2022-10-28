package main

import (
	"context"
	"encoding/json"
	"input-harga-storage-service/config"
	"log"
)

type HargaRequest struct {
	ID           string `json:"id"`
	HargaTopup   int    `json:"harga_topup"`
	HargaBuyback int    `json:"harga_buyback"`
	AdminID      string `json:"admin_id"`
}

func ReadMessage() {
	c := config.NewConfig()
	ctx := context.Background()
	for {
		m, err := c.Kr.FetchMessage(ctx)
		if err != nil {
			break
		}
		log.Printf("message at topic/partition/offset %v/%v/%v: %s\n", m.Topic, m.Partition, m.Offset, string(m.Key))
		var hr HargaRequest
		if err := json.Unmarshal(m.Value, &hr); err != nil {
			log.Printf("unmarshall data error : %s", err.Error())
		}

		res, err := c.Db.Exec(
			"UPDATE harga SET id=$1, harga_topup=$2, harga_buyback=$3, created_by=$4",
			hr.ID, hr.HargaTopup, hr.HargaBuyback, hr.AdminID)

		if err != nil {
			log.Printf("Update Failed error : %s", err.Error())
		}
		if rowsAffected, _ := res.RowsAffected(); rowsAffected == 0 {
			err = c.Db.QueryRow(
				"INSERT INTO harga(id, harga_topup, harga_buyback, created_by) VALUES($1, $2, $3, $4)",
				hr.ID, hr.HargaTopup, hr.HargaBuyback, hr.AdminID).Err()

			if err != nil {
				log.Printf("Create Failed error : %s", err.Error())
			}
		}

		if err := c.Kr.CommitMessages(ctx, m); err != nil {
			log.Printf("CommitMessage Failed error : %s", err.Error())
		}

	}
}
func main() {
	ReadMessage()
}
