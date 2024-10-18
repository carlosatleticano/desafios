package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
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
	http.ListenAndServe(":8070",nil)
}
	
func BuscaCotacaoHandler (w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/cotacao" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		cotacao, msg , erro := BuscaCotacao()

		if erro != nil {
			w.WriteHeader(http.StatusInternalServerError)
			resposta := map[string]string{"Erro":msg}
			json.NewEncoder(w).Encode(resposta)
			return	
		}
		w.Header().Set("Conten-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		valor_conv, _ := strconv.ParseFloat(cotacao.USDBRL.Bid, 64)
		resposta := map[string]float64{"Valor":valor_conv}
		json.NewEncoder(w).Encode(resposta)
}

func BuscaCotacao() (*CotacaoMoeada, string, error) {
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, 200 * time.Millisecond)  
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
		if err != nil {	
					println("Erro ao fazer a requisição: ")
					return nil, "Contexto da requisicao excedido", err
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {			
			println("Erro ao fazer a requisição: ")
			return nil, "Contexto da requisicao excedido", err
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {			panic(err)		}
		var cotacao CotacaoMoeada
		err = json.Unmarshal(body, &cotacao)
		if err != nil {			panic(err)		}
		
		msg, err  := RegistraCotacao(cotacao.USDBRL.Bid)
		if err	!= nil {
			return nil, msg, err
		}
		return &cotacao, "",  nil
}

func RegistraCotacao(cotacaoAtual string) (msg string, err error) {
	db, err := gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Falha ao abrir o banco de dados: ", err)
	}
	
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10 * time.Millisecond)  
	defer cancel()

    db.AutoMigrate(&Cotacao{})

	valorConv, _ := strconv.ParseFloat(cotacaoAtual, 64)
	cotacao := Cotacao{Valor: valorConv}
	err = db.WithContext(ctx).Create(&cotacao).Error
	if err != nil {
		return "Contexto do banco excedido", err
	}
	println("Cotação registrada com sucesso")
	return "", nil
}
