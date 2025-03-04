package main

import (
	"fmt"
	"log"
    "os"
    "github.com/joho/godotenv"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tinchomorilla/light-wallet/account"
)

func main() {
    
    api_key := get_api_key()
    
    // Connect to the Ethereum network
    client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + api_key)

    if err != nil {
        log.Fatalf("Error al conectar con Infura: %v", err)
    }

    // Generate a new Ethereum address
    address, err := account.GenerateNewAccount("hola") 
    if err != nil {
        log.Fatalf("Error generating new account: %v", err)
    }

    // Get the balance of the new address
    balance, err := account.GetBalance(client, address)
    if err != nil {
        log.Fatalf("Error fetching balance: %v", err)
    }

    fmt.Printf("Balance of address %s: %s ETH\n", address.Hex(), balance.String())

    fmt.Println("Generated new Ethereum address:", address.Hex())

    fmt.Println("Conexión exitosa a Infura:", client)
}

func get_api_key() (string) {
    // Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the API key from the environment variable
	apiKey := os.Getenv("API_KEY_INFURA")
    return apiKey
}