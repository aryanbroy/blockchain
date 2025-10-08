package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
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

func GenratePrivateKey() ecdsa.PrivateKey {
	pvtKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		log.Panicln("Error generating private key")
	}
	return *pvtKey
}

func (wa *Wallet) GenerateAddress() {
	wa.PrivateKey = GenratePrivateKey()

	ecdsaPubKey := &wa.PrivateKey.PublicKey
	pubKey, err := x509.MarshalPKIXPublicKey(ecdsaPubKey)
	if err != nil {
		log.Panicln("Error marshaling public key: ", err)
	}
	wa.PublicKey = pubKey

	shaHasher := sha256.New()
	_, err = shaHasher.Write(pubKey)
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

	fmt.Printf("Original: %x\nEndoced: %s\n", binaryAddr, base58Addr)
}
