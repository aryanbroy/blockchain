package main

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"

	// "crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
	"slices"
)

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

const subsidy = 10

func NewCoinbaseTx(to, data string) *Transaction {
	txIn := TXInput{[]byte{}, -1, []byte(data), nil}
	// txIn := TXInput{[]byte{}, -1, data}
	txOut := TXOutput{subsidy, to}

	transaction := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}
	return &transaction
}

func (bc *Blockchain) FindUnspentTx(address string) []Transaction {
	spentTx := make(map[string][]int)
	var unspentTx []Transaction

	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			txId := hex.EncodeToString(tx.ID)
			hasUnspent := false

			if !tx.isCoinBase() {
				for _, in := range tx.Vin {
					// check: if this output can be unlocked by the user's pubkey
					inId := hex.EncodeToString(in.Txid)
					spentTx[inId] = append(spentTx[inId], in.Vout)
				}
			}

			for voutIdx, out := range tx.Vout {
				isSpent := slices.Contains(spentTx[txId], voutIdx)

				if !isSpent && out.CanBeUnlockedWith(address) {
					hasUnspent = true
				}
			}
			if hasUnspent {
				unspentTx = append(unspentTx, *tx)
			}
		}

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}

	return unspentTx
}

func NewUTXOTransaction(from, to string, amount int, bc *Blockchain) *Transaction {
	var txInputs []TXInput
	var txOutputs []TXOutput

	acc, validOutputs := bc.FindSpendableOutputs(from, amount)
	if acc < amount {
		log.Panicln("Error! Not enough funds")
	}

	for txId, outputs := range validOutputs {
		txid, err := hex.DecodeString(txId)
		if err != nil {
			log.Panicln("error decoding transaction id")
		}

		for _, output := range outputs {
			currInput := TXInput{
				Txid: txid,
				Vout: output,
				// ScriptSig: from,

			}
			txInputs = append(txInputs, currInput)
		}
	}

	currOutput := TXOutput{
		Value:        amount,
		ScriptPubKey: to,
	}
	txOutputs = append(txOutputs, currOutput)

	if acc > amount {
		changeOutput := TXOutput{
			Value:        acc - amount,
			ScriptPubKey: from,
		}
		txOutputs = append(txOutputs, changeOutput)
	}

	tx := Transaction{
		ID:   nil,
		Vin:  txInputs,
		Vout: txOutputs,
	}

	tx.SetId()
	return &tx
}

func (tx *Transaction) HashTransaction() []byte {
	var buf bytes.Buffer

	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(tx)
	if err != nil {
		log.Panicln("Error encoding transaction: ", err)
	}

	txHash := sha256.Sum256(buf.Bytes())
	return txHash[:]
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
