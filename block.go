package main

import (
	"bytes"
	"encoding/gob"
)

// transforms incoming data to binary data to be transmitted further
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	encoder.Encode(b)

	return result.Bytes()
}

// does the opposite of Serialize
func Deserialize(data []byte) *Block {
	var block Block

	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	decoder.Decode(&block)

	return &block
}
