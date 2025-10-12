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

	// wallets.Wallets = make(map[string]*Wallet)
	// wallets.SaveFile()
	// fmt.Println("Wallets: ", wallets)

	wallets.LoadFromFile()

	return &wallets
}

func (ws *Wallets) SaveFile() {
	dummyWallet := NewWallet()
	dummyAddress := string(dummyWallet.GenerateAddress())

	privKey := dummyWallet.PrivateKey
	derBytes, err := x509.MarshalECPrivateKey(&privKey)
	if err != nil {
		log.Panicln("Error marshaling private key, x509 error: ", err)
	}

	ws.Wallets[dummyAddress] = dummyWallet

	storeWallet := &StorageWallet{}
	storeWallet.Curve = "P-256"
	storeWallet.PrivateDer = derBytes
	storeWallet.PubKeyByte = dummyWallet.PublicKey
	storeWallet.Address = dummyAddress

	jsonData, err := json.Marshal(storeWallet)

	// jsonData, err := json.Marshal(ws)
	if err != nil {
		log.Panicln("Error marshaling data: ", err)
	}

	err = os.WriteFile(walletFile, jsonData, 0544)
	if err != nil {
		log.Panicln("Error writing wallet data to file: ", err)
	}
}

func (ws *Wallets) LoadFromFile() {
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

	fmt.Println(storeWallet)
}
