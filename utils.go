package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
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
