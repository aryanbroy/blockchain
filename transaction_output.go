package main

import "encoding/hex"

type TXOutput struct {
	Value        int
	ScriptPubKey string
}

func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
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
