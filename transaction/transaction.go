package transaction

import (
	"context"
	"fmt"
	"math/big"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendTransaction dynamically fetches nonce, gas price, and gas limit
func SendTransaction(client *ethclient.Client, privateKey *bind.TransactOpts, toAddress common.Address, value *big.Int) (common.Hash, error) {
	ctx := context.Background()

	// Get the nonce (transaction count of sender)
	nonce, err := client.PendingNonceAt(ctx, privateKey.From)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get current gas price from the network
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to get gas price: %v", err)
	}

	// Estimate gas limit based on transaction type
	msg := ethereum.CallMsg{
		From:  privateKey.From,
		To:    &toAddress,
		Value: value,
		Data:  nil, // No data means it's a simple ETH transfer
	}
	gasLimit, err := client.EstimateGas(ctx, msg)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to estimate gas limit: %v", err)
	}

	// Create the transaction
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	// Sign the transaction
	signedTx, err := privateKey.Signer(privateKey.From, tx)
	if err != nil {
		return common.Hash{},fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send the transaction
	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to send transaction: %v", err)
	}

	// Return the transaction hash to the user 
	// so they can track the transaction
	return signedTx.Hash(), nil
}
