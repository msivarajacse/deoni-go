package block

import (
	"crypto/ed25519"
	"crypto/sha256"
	"deoni/utils"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// TODO:/* Block Creation Code */

type Block struct {
	no           int64
	previousHash [32]byte
	timestamp    int64
	transactions []*Transaction
}

func NewBlock(no int64, previousHash [32]byte, transactions []*Transaction) *Block {
	block := new(Block)
	block.timestamp = time.Now().UnixNano()
	block.previousHash = previousHash
	block.no = no
	block.transactions = transactions
	return block
}

func (block *Block) Print() {
	fmt.Printf("No %d\n", block.no)
	fmt.Printf("Timestamp %d\n", block.timestamp)
	fmt.Printf("Previous_hash %x\n", block.previousHash)
	for _, t := range block.transactions {
		t.Print()
	}
}

func (block *Block) Hash() [32]byte {
	data, _ := json.Marshal(block)
	fmt.Println(string(data))
	return sha256.Sum256([]byte(data))
}

func (block *Block) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Timestamp    int64          `json:"timestamp"`
		No           int            `json:"no"`
		PreviousHash [32]byte       `json:"previous_hash"`
		Transactions []*Transaction `json:"transactions"`
	}{
		Timestamp:    block.timestamp,
		No:           int(block.no),
		PreviousHash: block.previousHash,
		Transactions: block.transactions,
	})
}

// TODO:/* Blockchain Code */

type Blockchain struct {
	transactionPool []*Transaction
	chain           []*Block
}

func NewBlockchain() *Blockchain {
	block := &Block{}
	blockchain := new(Blockchain)
	blockchain.CreateBlock(1, block.Hash())
	return blockchain
}

func (blockchain *Blockchain) CreateBlock(no int64, previousHash [32]byte) *Block {
	block := NewBlock(no, previousHash, blockchain.transactionPool)
	blockchain.chain = append(blockchain.chain, block)
	blockchain.transactionPool = []*Transaction{}
	return block
}

func (block *Blockchain) LastBlock() *Block {
	return block.chain[len(block.chain)-1]
}

func (blockchain *Blockchain) Print() {
	for i, block := range blockchain.chain {
		fmt.Printf("%s Chain %d %s\n", strings.Repeat("=", 25), i, strings.Repeat("=", 25))
		block.Print()
	}
}

func (blockchain *Blockchain) AddTransaction(sender string, reciever string, amount float32, publicKey []byte, signature *utils.Signature) bool {
	transaction := NewTransaction(sender, reciever, amount)

	if blockchain.VerifyTransactionSignature(publicKey, signature, transaction) {
		//if blockchain.CalculateWalletBalance(sender) < amount {
		//	log.Println("ERROR: Sender Don't Have the Money")
		//	return false
		//}
		blockchain.transactionPool = append(blockchain.transactionPool, transaction)
		return true
	} else {
		fmt.Println("Error: Signature Verify Transaction")
	}
	return false
}

func (blockchain *Blockchain) VerifyTransactionSignature(publicKey []byte, signature *utils.Signature, transaction *Transaction) bool {
	data, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(data))
	return ed25519.Verify(publicKey, hash[:], signature.Sign)
}

func (blockchain *Blockchain) CopyTransactionPool() []*Transaction {
	transactions := make([]*Transaction, 0)
	for _, transaction := range blockchain.transactionPool {
		transactions = append(transactions, NewTransaction(transaction.sender, transaction.reciever, transaction.amount))
	}
	return transactions
}

func (blockchain *Blockchain) CalculateWalletBalance(address string) float32 {
	var balance float32 = 0
	for _, block := range blockchain.chain {
		for _, transaction := range block.transactions {
			amount := transaction.amount
			if address == transaction.reciever {
				balance += amount
			}
			if address == transaction.sender {
				balance -= amount
			}
		}
	}
	return balance
}

type Transaction struct {
	sender   string
	reciever string
	amount   float32
}

func NewTransaction(sender string, reciever string, amount float32) *Transaction {
	return &Transaction{sender, reciever, amount}
}

func (transaction *Transaction) Print() {
	fmt.Printf("%s\n", strings.Repeat("-", 40))
	fmt.Printf(" Sender     %s\n", transaction.sender)
	fmt.Printf(" Reciever   %s\n", transaction.reciever)
	fmt.Printf(" Amount     %f\n", transaction.amount)
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender"`
		Reciever string  `json:"reciever"`
		Amount   float32 `json:"amount"`
	}{
		Sender:   transaction.sender,
		Reciever: transaction.reciever,
		Amount:   transaction.amount,
	})
}
