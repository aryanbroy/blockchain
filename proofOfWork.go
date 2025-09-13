package main

import (
	"bytes"
	"fmt"
	"math/big"
)

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256-targetBits)

	pow := &ProofOfWork{
		block:  b,
		target: target,
	}
	return pow
}

func IntToHex(data int64) []byte {
	hexString := fmt.Sprintf("%x", data)
	return []byte(hexString)
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join([][]byte{
		pow.block.PreviousBlockHash,
		pow.block.Data,
		IntToHex(pow.block.Timestamp),
		IntToHex(int64(targetBits)),
		IntToHex(int64(nonce)),
	}, []byte{})

	return data
}
