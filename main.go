package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"errors"	
)

type Address struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
}

type ViaCEPResponse struct {
	Cep        string `json:"cep"`
	Logradouro string `json:"logradouro"`
	Bairro     string `json:"bairro"`
	Localidade string `json:"localidade"`
	Uf         string `json:"uf"`
}

type BrasilAPIResponse struct {
	Cep        string `json:"cep"`
	Street     string `json:"street"`
	Neighborhood string `json:"neighborhood"`
	City       string `json:"city"`
	State      string `json:"state"`
}

type APIResponse struct {
	API      string
	Address  Address
	Error    error
}

func fetchFromAPI(ctx context.Context, url string, apiName string, ch chan<- APIResponse) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		ch <- APIResponse{API: apiName, Error: err}
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- APIResponse{API: apiName, Error: err}
		return
	}
	defer resp.Body.Close()

	var address Address
	if apiName == "BrasilAPI" {
		var brasilAPI BrasilAPIResponse
		if err := json.NewDecoder(resp.Body).Decode(&brasilAPI); err != nil {
			ch <- APIResponse{API: apiName, Error: err}
			return
		}

		if brasilAPI.Street == "" && brasilAPI.Neighborhood == "" && brasilAPI.City == "" {
			ch <- APIResponse{API: apiName, Error: errors.New("CEP não encontrado")}
			return
		}

		address = Address{
			Cep:        brasilAPI.Cep,
			Logradouro: brasilAPI.Street,
			Bairro:     brasilAPI.Neighborhood,
			Localidade: brasilAPI.City,
			Uf:         brasilAPI.State,
		}

	} else if apiName == "ViaCEP" {
		var viaCEP ViaCEPResponse
		if err := json.NewDecoder(resp.Body).Decode(&viaCEP); err != nil {
			ch <- APIResponse{API: apiName, Error: err}
			return
		}

		if viaCEP.Logradouro == "" && viaCEP.Bairro == "" && viaCEP.Localidade == "" {
			ch <- APIResponse{API: apiName, Error: errors.New("CEP não encontrado")}
			return
		}
		
		address = Address{
			Cep:        viaCEP.Cep,
			Logradouro: viaCEP.Logradouro,
			Bairro:     viaCEP.Bairro,
			Localidade: viaCEP.Localidade,
			Uf:         viaCEP.Uf,
		}
	}

	ch <- APIResponse{API: apiName, Address: address}
}

func main() {
	var cep string
	fmt.Print("Digite o CEP (somente números): ")
	fmt.Scanln(&cep)

	timeout := time.Second

	api1 := "https://brasilapi.com.br/api/cep/v1/" + cep
	api2 := "http://viacep.com.br/ws/" + cep + "/json/"

	responseChannel := make(chan APIResponse)

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	go fetchFromAPI(ctx, api1, "BrasilAPI", responseChannel)
	go fetchFromAPI(ctx, api2, "ViaCEP", responseChannel)

	select {
	case result := <-responseChannel:
		if result.Error != nil {
			fmt.Println("Erro ao buscar dados da API:", result.API, "-", result.Error)
		} else {
			fmt.Printf("API mais rápida: %s\n", result.API)
			fmt.Printf("Endereço: CEP: %s, Logradouro: %s, Bairro: %s, Cidade: %s, UF: %s\n",
				result.Address.Cep, result.Address.Logradouro, result.Address.Bairro, result.Address.Localidade, result.Address.Uf)
		}
	case <-ctx.Done():
		fmt.Println("Timeout: Nenhuma API respondeu dentro de 1 segundo.")
	}
}
