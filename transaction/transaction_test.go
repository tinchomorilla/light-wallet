package transaction_test

import (
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/tinchomorilla/light-wallet/transaction" // Your package where SendTransaction is located
)

// Test the sendTransaction function to check
// if it sends a transaction correctly
// This test sends 0.0005 ETH from one account to another
func TestSendTransaction(t *testing.T) {

	err := godotenv.Load("/home/tincho/Documents/light-wallet/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	api_key := os.Getenv("API_KEY_INFURA")

	// Connect to the Ethereum network
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + api_key)

	if err != nil {
		
		log.Fatalf("Error connecting to infura %v", err)
	}


	addressOnePrivateKey := os.Getenv("ACCOUNT_ONE_PRIVATE_KEY")

	addressDestination := os.Getenv("ACCOUNT_TWO_ADDRESS")

	// Convert private key to *ecdsa.PrivateKey
	privateKeyECDSA, err := crypto.HexToECDSA(addressOnePrivateKey)
    if err != nil {
        log.Fatalf("Error converting private key: %v", err)
    }

	// Convert private key to TransactOpts for signing the transaction
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKeyECDSA, big.NewInt(11155111)) // Chain ID 11155111 is for Ropsten
	if err != nil {
		t.Fatalf("Failed to create transact options: %v", err)
	}

	// Convert address to common.Address 
	toAddress := common.HexToAddress(addressDestination)

	// Value to send
	value := big.NewInt(500000000000000)  // 0.0005 ETH in wei


	// Call SendTransaction
	_, err = transaction.SendTransaction(client, transactOpts, toAddress, value)
	if err != nil {
		t.Fatalf("Failed to send transaction: %v", err)
	}

	// Here we can check the transaction hash and status
	// to ensure the transaction was successful
	// I checked it manually in my wallets and it worked :)
}
