package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Dolar struct {
	Dolar string `json:"dolar"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao cria um novo request: %v\n", err)
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao fazer a requisição: %v\n", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusInternalServerError {
		fmt.Println("Erro ao solicitar a cotação")
		return
	}

	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao ler a requisição: %v\n", err)
		return
	}

	var dolar Dolar
	dolar.Dolar = string(resbody)
	file, err := os.Create("cotacao.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao criar o arquivo: %v\n", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("Dólar:{%s}", dolar.Dolar))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao persistir o valor no arquivo: %v\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Cotação salva com sucesso: Dólar:{%s}\n", dolar.Dolar)
}
