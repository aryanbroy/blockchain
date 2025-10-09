package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
)

type Wallets struct {
	Wallets map[string]*Wallet
}

func NewWallets() *Wallets {
	wallets := Wallets{}

	dummyWallet := NewWallet()
	dummyAddress := dummyWallet.GenerateAddress()

	wallets.Wallets = make(map[string]*Wallet)
	wallets.Wallets[string(dummyAddress)] = dummyWallet

	return &wallets
}

func (wa *Wallets) writeToFile() {
	var data bytes.Buffer

	encoder := gob.NewEncoder(&data)
	err := encoder.Encode(wa)
	if err != nil {
		log.Panicln("Error encoding wallets")
	}

	err = os.WriteFile(walletFile, data.Bytes(), 0777)
	if err != nil {
		log.Panicln("Error writing wallets data to a file")
	}
}
