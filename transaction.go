package main

import "fmt"

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

const subsidy = 1

func NewCoinbaseTx(to, data string) *Transaction {
	if data == "" {
		data = fmt.Sprintf("Reward to '%s'", to)
	}

	txIn := TXInput{[]byte{}, -1, data}
	txOut := TXOutput{subsidy, to}
	transaction := Transaction{nil, []TXInput{txIn}, []TXOutput{txOut}}

	return &transaction
}
