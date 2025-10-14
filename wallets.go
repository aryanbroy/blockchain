package main

import (
	"crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet `json:"wallets"`
}

func NewWallets() *Wallets {
	wallets := Wallets{}

	wallets.Wallets = make(map[string]*Wallet)
	wallets.LoadFromFile()

	return &wallets
}

func (ws *Wallets) LoadFromFile() {
	_, err := os.Stat(walletFile)
	if os.IsNotExist(err) {
		return
	}
	byteData, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panicln("Error reading from wallet file")
	}

	fmt.Println(string(byteData))

	storeWallet := &StorageWallet{}
	err = json.Unmarshal(byteData, storeWallet)
	if err != nil {
		log.Panicln("Error unmarshaling byte data to store wallet")
	}

	address := storeWallet.Address

	prvKey, err := x509.ParseECPrivateKey(storeWallet.PrivateDer)
	if err != nil {
		log.Panicln("Error parsing der bytes to private key: ", err)
	}

	wallet := &Wallet{
		PrivateKey: *prvKey,
		PublicKey:  storeWallet.PubKeyByte,
	}

	wallets := make(map[string]*Wallet)
	wallets[address] = wallet

	ws.Wallets = wallets
}

func (ws *Wallets) CreateWallet() string {
	wallet := NewWallet()
	address := string(wallet.GenerateAddress())

	ws.Wallets[address] = wallet

	privKey := wallet.PrivateKey
	derBytes, err := x509.MarshalECPrivateKey(&privKey)
	if err != nil {
		log.Panicln("Error marshaling private key, x509 error: ", err)
	}

	ws.Wallets[address] = wallet

	storeWallet := &StorageWallet{}
	storeWallet.Curve = "P-256"
	storeWallet.PrivateDer = derBytes
	storeWallet.PubKeyByte = wallet.PublicKey
	storeWallet.Address = address

	jsonData, err := json.Marshal(storeWallet)
	if err != nil {
		log.Panicln("Error marshaling data: ", err)
	}

	err = os.WriteFile(walletFile, jsonData, 0544)
	if err != nil {
		log.Panicln("Error writing wallet data to file: ", err)
	}

	return address
}
