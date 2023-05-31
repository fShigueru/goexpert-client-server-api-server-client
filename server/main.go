package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	_ "github.com/mattn/go-sqlite3"
)

type RealDolarResponse struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type RealDolar struct {
	ID         int    `gorm:"primaryKey"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
	gorm.Model
}

func main() {
	println("Server :8080")
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	realDolar, err := BuscaCambioRealDolar()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(realDolar.Usdbrl.Bid)
}

func BuscaCambioRealDolar() (*RealDolarResponse, error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var realDolarResponse RealDolarResponse
	err = json.Unmarshal(body, &realDolarResponse)
	if err != nil {
		return nil, err
	}

	Sqlite(&realDolarResponse)
	return &realDolarResponse, nil
}

func Sqlite(realDolarResponse *RealDolarResponse) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := gorm.Open(sqlite.Open("db/realdolar.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	println("Conectando sqlite3")
	db.AutoMigrate(&RealDolar{})

	usdbrl := realDolarResponse.Usdbrl
	db.WithContext(ctx).Create(&RealDolar{
		Code:       usdbrl.Code,
		Codein:     usdbrl.Codein,
		Name:       usdbrl.Name,
		High:       usdbrl.High,
		Low:        usdbrl.Low,
		VarBid:     usdbrl.VarBid,
		PctChange:  usdbrl.PctChange,
		Bid:        usdbrl.Bid,
		Ask:        usdbrl.Ask,
		Timestamp:  usdbrl.Timestamp,
		CreateDate: usdbrl.CreateDate,
	})

	var realDolarBusca []RealDolar
	db.Find(&realDolarBusca)
	for _, realDolar := range realDolarBusca {
		fmt.Println(realDolar.Name, realDolar.Bid)
	}
}
