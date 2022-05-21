package wallet

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"deoni/utils"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
)

type Wallet struct {
	privateKey []byte
	publicKey  []byte
	pubaddress string
}

func NewWallet() *Wallet {
	wallet := new(Wallet)
	publicKey, privateKey, _ := ed25519.GenerateKey(rand.Reader)
	wallet.privateKey = privateKey
	wallet.publicKey = publicKey
	//fmt.Printf("%T", wallet.publicKey)
	h2 := sha256.New()
	h2.Write(wallet.publicKey)
	digest2 := h2.Sum(nil)
	h3 := ripemd160.New()
	h3.Write(digest2)
	digest3 := h3.Sum(nil)
	vd4 := make([]byte, 21)
	vd4[0] = 0x00
	copy(vd4[1:], digest3[:])
	h5 := sha256.New()
	h5.Write(vd4)
	digest5 := h5.Sum(nil)
	h6 := sha256.New()
	h6.Write(digest5)
	digest6 := h6.Sum(nil)
	chsum := digest6[:4]
	dc8 := make([]byte, 25)
	copy(dc8[:21], vd4[:])
	copy(dc8[21:], chsum[:])

	address := base58.Encode(dc8)
	wallet.pubaddress = address
	//fmt.Println(address)
	return wallet
}

func (wallet *Wallet) PrivateKey() []byte {
	return wallet.privateKey
}

func (wallet *Wallet) PrivateKeyHex() string {
	return hex.EncodeToString(wallet.privateKey)
}

func (wallet *Wallet) PublicKey() []byte {
	return wallet.publicKey
}

func (wallet *Wallet) PublicKeyHex() string {
	return hex.EncodeToString(wallet.publicKey)
}

func (wallet *Wallet) PublicAddress() string {
	return wallet.pubaddress
}

type Transaction struct {
	senderPrivateKey []byte
	senderPublicKey  []byte
	senderAddress    string
	recieverAddress  string
	amount           float32
}

func NewTransaction(privateKey []byte, publicKey []byte, sender string, reciever string, amount float32) *Transaction {
	return &Transaction{privateKey, publicKey, sender, reciever, amount}
}

func (transaction *Transaction) GenerateSignature() *utils.Signature {
	data, _ := json.Marshal(transaction)
	hash := sha256.Sum256([]byte(data))
	sign := ed25519.Sign(transaction.senderPrivateKey, hash[:])
	return &utils.Signature{Sign: sign}
}

func (transaction *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Sender   string  `json:"sender"`
		Reciever string  `json:"reciever"`
		Amount   float32 `json:"amount"`
	}{
		Sender:   transaction.senderAddress,
		Reciever: transaction.recieverAddress,
		Amount:   transaction.amount,
	})
}
