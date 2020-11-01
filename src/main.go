package src

import (
	"encoding/json"
	"log"
	"net/http"
)

type Body struct{
	Mail string
	Link string
}

func (db *DataBaseWork)Subscribe(w http.ResponseWriter, r *http.Request) {
	log.Println("In Subscribe")
	if r.Method != http.MethodPost {
		log.Println("Incorrect method")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	body := new(Body)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := db.addRecord(body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}