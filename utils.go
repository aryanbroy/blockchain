package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func (tx *Transaction) SetId() {
	var encodedId bytes.Buffer
	var txId [32]byte

	encoder := gob.NewEncoder(&encodedId)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln("Error encoding the transaction")
	}

	txId = sha256.Sum256(encodedId.Bytes())
	tx.ID = txId[:]
}

func (tx Transaction) isCoinBase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func WriteToJson(ws *Wallets, sws *StorageWallets) string {
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
	if err != nil {
		return nil, err
	}

	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return jsonData, err
}

func Sign(tx *Transaction, prvKey *ecdsa.PrivateKey) []byte {
	txHash := tx.HashTransaction()
	fmt.Println(string(txHash))

	signBytes, err := ecdsa.SignASN1(rand.Reader, prvKey, txHash)
	if err != nil {
		log.Panicln("Error signing a signature: ", err)
	}

	return signBytes
}

func Verify(tx *Transaction) bool {
	isVerified := false
	return isVerified
}
