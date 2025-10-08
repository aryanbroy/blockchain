package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/btcsuite/btcutil/base58"
	"github.com/btcsuite/golangcrypto/ripemd160"
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type Wallets struct {
	Wallets map[string]*Wallet
}

const version = 0

func NewWallet() *Wallet {
	prvKey, pubKey := GenerateKeyPair()

	return &Wallet{
		PrivateKey: prvKey,
		PublicKey:  pubKey,
	}
}

func GenerateKeyPair() (ecdsa.PrivateKey, []byte) {
	var pubKey []byte

	ec := elliptic.P256()
	prvKey, err := ecdsa.GenerateKey(ec, rand.Reader)
	if err != nil {
		log.Panicln("Error generating private key")
	}

	pubKey = append(prvKey.X.Bytes(), prvKey.Y.Bytes()...)
	return *prvKey, pubKey
}

func (wa *Wallet) GenerateAddress() {
	shaHasher := sha256.New()
	_, err := shaHasher.Write(wa.PublicKey)
	if err != nil {
		log.Panicln("Error writing in sha function: ", err)
	}
	data := shaHasher.Sum(nil)

	ripemdHasher := ripemd160.New()
	_, err = ripemdHasher.Write(data)
	if err != nil {
		log.Panicln("Error writing in ripemd function: ", err)
	}
	ripemdResult := ripemdHasher.Sum(nil)

	versionPayload := append([]byte{version}, ripemdResult...)

	hash := sha256.Sum256(versionPayload)
	resultHash := sha256.Sum256(hash[:])
	addrCheckSum := resultHash[:4]

	binaryAddr := append(versionPayload, addrCheckSum...)
	base58Addr := base58.Encode(binaryAddr)

	fmt.Printf("Original: %x\nEncoded: %s\n", binaryAddr, base58Addr)
}
