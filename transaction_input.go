package main

type TXInput struct {
	Txid []byte
	Vout int
	// ScriptSig string
	PubKey    []byte
	Signature []byte
}

// func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
// 	return in.ScriptSig == unlockingData
// }
