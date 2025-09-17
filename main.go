package main

import (
	"math"
	"math/big"
	"os"
	"time"
)

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const targetBits = 24
const maxNonce = math.MaxInt64

func main() {
	bc, err := NewBlockChain()
	if err != nil {
		os.Exit(0)
	}

	defer bc.db.Close()

	cli := CLI{bc: bc}
	cli.Run()
}

func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock() *Block {
	return NewBlock("Genesis block", []byte{})
}
