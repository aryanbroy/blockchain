package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"time"
)

type Block struct {
	Timestamp         int64
	Data              []byte
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

type BlockChain struct {
	blocks []*Block
}

type ProofOfWork struct {
	block  *Block
	target *big.Int
}

const targetBits = 24
const maxNonce = math.MaxInt64

func main() {
	bc := NewBlockChain()

	bc.AddBlock("Send 1 ETH to alom")
	bc.AddBlock("Send 1.5 ETH to 6wong")
	bc.AddBlock("Send 2 ETH to achi hoon")

	for _, block := range bc.blocks {
		pow := NewProofOfWork(block)
		isValid := pow.Validate()
		fmt.Printf("Previous hash: %x\n", block.PreviousBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Pow: %v\n", isValid)
		fmt.Println()
	}
}

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

func (pow *ProofOfWork) Run() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%s\"\n", pow.block.Data)
	for nonce < maxNonce {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x", hash)
		hashInt.SetBytes(hash[:])

		if nonce == math.MaxInt64 {
			fmt.Println("Max nonce reached")
			break
		}
		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++
		}
	}

	fmt.Print("\n\n")
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
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

func NewBlockChain() *BlockChain {
	return &BlockChain{[]*Block{NewGenesisBlock()}}
}

func (bc *BlockChain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}
