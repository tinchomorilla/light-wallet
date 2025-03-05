package account_test

import (
	"fmt"
	"log"
	"math/big"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/tinchomorilla/light-wallet/account"
)


// The purpose of this test is to check if the balance 
// of an Ethereum address is correct, in particular my address
// which has 0.03 SepoliaETH.
func TestGetBalance(t *testing.T) {

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

	// My metamask address
	address := os.Getenv("ACCOUNT_ONE_ADDRESS")

	// Convert address to common.Address
	commonAddress := common.HexToAddress(address)

	// Get the balance of my address, who has 0.03 SepoliaETH
	balance, err := account.GetBalance(client, commonAddress)
	if err != nil {
		log.Fatalf("Error fetching balance: %v", err)
	}

	// Convert to ETH
	balanceInETH := new(big.Float).SetInt(balance)
	balanceInETH.Quo(balanceInETH, big.NewFloat(1e18))

	// Round the balance to 2 decimal places
	roundedBalance := new(big.Float).SetPrec(64).SetMode(big.ToNearestEven).Quo(balanceInETH, big.NewFloat(0.01)) // Divide by 0.01 to round
	roundedBalance.Mul(roundedBalance, big.NewFloat(0.01))  // Multiply by 0.01 to get the rounded value back

	// Convert to string with 2 decimal places
	roundedBalanceStr := fmt.Sprintf("%.2f", roundedBalance)

	expectedBalance := "0.03"

	// Assert the rounded balance is correct
	assert.Equal(t, expectedBalance, roundedBalanceStr)
}
