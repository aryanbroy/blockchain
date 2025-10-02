package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

type Block struct {
	Timestamp         int64
	Transactions      []*Transaction
	PreviousBlockHash []byte
	Hash              []byte
	Nonce             int
}

// transforms incoming data to binary data to be transmitted further
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	if err := encoder.Encode(b); err != nil {
		log.Panicln("Error encoding buffer: ", err)
	}

	return result.Bytes()
}

// does the opposite of Serialize
func Deserialize(data []byte) *Block {
	var block Block

	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	if err := decoder.Decode(&block); err != nil {
		log.Panicln("Error decoding block: ", err)
	}
	return &block
}

func NewBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce
	return block
}

func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{})
}
