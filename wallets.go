package main

import (
	"crypto/x509"
	"encoding/json"
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
	if err != nil {
		log.Panicln("Error reading json: wallet file, ", err)
	}

	err = json.Unmarshal(data, sws)
	if err != nil {
		log.Panicln("Error unmarshaling into storageWallets: ", err)
	}

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

	sws.Wallets[address] = storeWallet

	jsonData, err := json.MarshalIndent(sws, "", "  ")
	if err != nil {
		log.Panicln("Error marshaling data: ", err)
	}

	err = os.WriteFile(walletFile, jsonData, 0777)
	if err != nil {
		log.Panicln("Error writing wallet json data to file: ", err)
	}

	return address
}

func ReadJson(filepath string) ([]byte, error) {
	_, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		log.Panicln("No such file exists, error reading json")
	}
	if err != nil {
		log.Panicln("Error reading json file")
	}

	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return jsonData, err
}
