package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GenerateNewAccount generates a new Ethereum account and returns the address, along with the encrypted private key.
func GenerateNewAccount(password string) (common.Address, *ecdsa.PrivateKey, error) {

	if len(password) == 0 {
		return common.Address{}, nil, fmt.Errorf("password cannot be empty")
	}

	// Create a new keypair
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return common.Address{}, nil, fmt.Errorf("failed to generate private key: %v", err)
	}

	// Derive the account address from the public key
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	return address, privateKey, nil
}

// GetBalance fetches the balance of an Ethereum address.
func GetBalance(client *ethclient.Client, address common.Address) (*big.Int, error) {
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get balance: %v", err)
	}
	return balance, nil
}
