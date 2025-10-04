package main

import (
	"encoding/hex"
	"fmt"
	"log"
	"slices"
)

type TXInput struct {
	Txid      []byte
	Vout      int
	ScriptSig string
}

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

const subsidy = 10

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TXInput{[]byte{}, -1, data}
	txOut := TXOutput{subsidy, to}
	transaction := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}

	return &transaction
}

func (tx Transaction) isCoinBase() bool {
	return len(tx.Vin) == 1 && tx.Vin[0].Vout == -1 && len(tx.Vin[0].Txid) == 0
}

func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
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
					if in.CanUnlockOutputWith(address) {
						inId := hex.EncodeToString(in.Txid)
						spentTx[inId] = append(spentTx[inId], in.Vout)
					}
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

func (bc *Blockchain) FindUTXOs(address string) []TXOutput {
	var utxos []TXOutput

	unspentTxs := bc.FindUnspentTx(address)
	for _, tx := range unspentTxs {
		for _, out := range tx.Vout {
			if out.CanBeUnlockedWith(address) {
				utxos = append(utxos, out)
			}
		}
	}

	return utxos
}

func (bc *Blockchain) FindSpendableOutputs(address string, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0
	unspentTx := bc.FindUnspentTx(address)

Outputs:
	for _, tx := range unspentTx {
		txId := hex.EncodeToString(tx.ID)

		for outputIdx, output := range tx.Vout {
			if output.CanBeUnlockedWith(address) {
				accumulated += output.Value
				unspentOutputs[txId] = append(unspentOutputs[txId], outputIdx)

				if accumulated >= amount {
					break Outputs
				}
			}
		}
	}

	return accumulated, unspentOutputs
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
				Txid:      txid,
				Vout:      output,
				ScriptSig: from,
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
