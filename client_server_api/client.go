package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Response struct {
    Valor float64 `json:"Valor"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300 * time.Millisecond)
	defer cancel()

    requisição, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8070/cotacao", nil)
    if err != nil {        log.Fatalf("Erro ao fazer a requisição: %v", err)    }
    resposta, err := http.DefaultClient.Do(requisição)
    if err != nil {       
		 println("Erro: Contexto da solicitacao ao servidor excedido:")  
		 return  
	}
    defer resposta.Body.Close()
    corpo, err := io.ReadAll(resposta.Body)
    if err != nil {        log.Fatalf("Erro ao ler o corpo da resposta: %v", err)    }

    var dados Response

	fmt.Println(string(corpo))

    err = json.Unmarshal(corpo, &dados)
    if err != nil {        log.Fatalf("Erro ao fazer o parse da resposta: %v", err)    }

    str2 := strconv.FormatFloat(dados.Valor, 'f', 6, 64)

    file , err := os.Create("cotacao.txt")
	if err != nil {	fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)   }
	defer file.Close()
	_, err = file.WriteString(fmt.Sprintf("Dólar: %s", str2))
	
}



