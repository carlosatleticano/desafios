package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type CepBrasilApi struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CepViacep struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Unidade     string `json:"unidade"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Estado      string `json:"estado"`
	Regiao      string `json:"regiao"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

func main() {
    c1 := make(chan string)
    c2 := make(chan string)
    for _, cep := range os.Args[1:] {  // recuoera o parametro ao chamar o pgm ex: go run main.go Olá
      go func () {
        //time.Sleep(2000 * time.Millisecond)
        req, err := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)}
        defer req.Body.Close()
        res, err := io.ReadAll(req.Body)
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao ler a resposta requisição: %v\n", res)           }
        var data CepViacep
        err = json.Unmarshal(res, &data)
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)           }
        c1 <- string(res)
      }()

      go func () {
        //time.Sleep(200 * time.Millisecond)
        req, err := http.Get("https://brasilapi.com.br/api/cep/v1/" + cep )
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao fazer requisição: %v\n", err)}
        defer req.Body.Close()
        res, err := io.ReadAll(req.Body)
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao ler a resposta requisição: %v\n", res)           }
        var data CepBrasilApi
        err = json.Unmarshal(res, &data)
        if err != nil {fmt.Fprintf(os.Stderr, "Erro ao fazer o parse da resposta: %v\n", err)           }
        c2 <- string(res)
      }()
      select {
        case msg1 := <-c1:
           fmt.Println("Recebido de ViaCep\n", msg1)
        case msg2 := <-c2:
           fmt.Println("Recebido de Brasilapi\n", msg2)
        case <-time.After(1000 * time.Millisecond):
           fmt.Println("Timeout")           
      }
    }
}



