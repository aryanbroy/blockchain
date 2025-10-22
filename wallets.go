package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

type StorageWallets struct {
	Wallets map[string]*StorageWallet `json:"wallets"`
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
	if err != nil {
		log.Panicln("Error describing a file: ", err)
	}

	byteData, err := os.ReadFile(walletFile)
	if err != nil {
		log.Panicln("Error reading from wallet file")
	}

	// fmt.Println(string(byteData))

	sws := &StorageWallets{}
	err = json.Unmarshal(byteData, sws)
	if err != nil {
		log.Panicln("Error unmarshaling byte data to store wallet")
	}

	wallets := sws.Wallets
	fmt.Println(wallets)
	// address := storeWallet.Address
	//
	// prvKey, err := x509.ParseECPrivateKey(storeWallet.PrivateDer)
	// if err != nil {
	// 	log.Panicln("Error parsing der bytes to private key: ", err)
	// }
	//
	// wallet := &Wallet{
	// 	PrivateKey: *prvKey,
	// 	PublicKey:  storeWallet.PubKeyByte,
	// }
	//
	// wallets := make(map[string]*Wallet)
	// wallets[address] = wallet
	//
	// ws.Wallets = wallets
}

func (ws *Wallets) CreateWallet() string {
	sws := &StorageWallets{Wallets: make(map[string]*StorageWallet)}

	data, err := ReadJson(walletFile)
	if os.IsNotExist(err) {
		address := WriteToJson(ws, sws)
		return address
	}
	if err != nil {
		log.Panicln("Error reading json: wallet file, ", err)
	}

	err = json.Unmarshal(data, sws)
	if err != nil {
		log.Panicln("Error unmarshaling into storageWallets: ", err)
	}

	address := WriteToJson(ws, sws)
	return address
}
