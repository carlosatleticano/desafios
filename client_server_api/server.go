package main

import (
	"encoding/json"
	"log"
	"io"
	"net/http"
	"context"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CotacaoMoeada struct {
	USDBRL struct {
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

type Cotacao struct {
	ID int `gorm:"primary_key"`
	Valor float64
	gorm.Model
}

func main() {
	http.HandleFunc("/cotacao", BuscaCotacaoHandler)
	http.ListenAndServe(":8080",nil)
}
	
func BuscaCotacaoHandler (w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cotacao" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cotacao, error := BuscaCotacao()
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return	
		}
		w.Header().Set("Conten-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		valor_conv, _ := strconv.ParseFloat(cotacao.USDBRL.Bid, 64)
		resposta := map[string]float64{"Valor":valor_conv}
		json.NewEncoder(w).Encode(resposta)
}

func BuscaCotacao() (*CotacaoMoeada, error) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 200 * time.Millisecond)  
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {			log.Fatal(err)		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {			log.Fatal(err)		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {			panic(err)		}
		var cotacao CotacaoMoeada
		err = json.Unmarshal(body, &cotacao)
		if err != nil {			panic(err)		}
		
		RegistraCotacao(cotacao.USDBRL.Bid)
		return &cotacao, nil
}

func RegistraCotacao(cotacaoAtual string) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Millisecond)  
	defer cancel()

	dsn := "root:root@tcp(localhost:3306)/goexpert?charset=utf8mb4&parseTime=True&loc=Local"
	db , err:= gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {panic("Falha ao se conectar ao banco de dados")}
    //db.Migrator().DropTable(&Cotacao{})
    db.AutoMigrate(&Cotacao{})

	valorConv, _ := strconv.ParseFloat(cotacaoAtual, 64)
	cotacao := Cotacao{Valor: valorConv}
	err = db.WithContext(ctx).Create(&cotacao).Error
	if err != nil {log.Fatalf("falha ao gravar o registro: %v", err)}
	println("Cotação registrada com sucesso")
}









