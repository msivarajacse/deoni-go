package main

import (
	"deoni/block"
	"deoni/wallet"
	"fmt"
)

func main() {
	walletA := wallet.NewWallet()
	walletB := wallet.NewWallet()
	//walletC := wallet.NewWallet()

	// Wallet Transaction Initiate
	transaction := wallet.NewTransaction(walletA.PrivateKey(), walletA.PublicKey(), walletA.PublicAddress(), walletB.PublicAddress(), 100)

	//	Create Block in Blockchain
	blockchain := block.NewBlockchain()
	isAdded := blockchain.AddTransaction(walletA.PublicAddress(), walletB.PublicAddress(), 100, walletA.PublicKey(), transaction.GenerateSignature())
	fmt.Println("Added?", isAdded)

	blockchain.CreateBlock(2, blockchain.LastBlock().Hash())

	fmt.Println("Balance of A is : ", blockchain.CalculateWalletBalance(walletA.PublicAddress()))
	fmt.Println("\n Balance of B is : ", blockchain.CalculateWalletBalance(walletB.PublicAddress()))
}
