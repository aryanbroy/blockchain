package main

import (
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type BlockChain struct {
	tip []byte
	db  *bolt.DB
}

const dbFile = "myDB.db"
const bucket = "bucket"

func NewBlockChain() (*BlockChain, error) {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Println("Error initializing database")
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			genesis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				log.Println("Error creating a bucket")
				return err
			}
			err = b.Put(genesis.Hash, genesis.Serialize())
			err = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
			if err != nil {
				log.Println("Error updating data in db")
				return err
			}
		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})

	return &BlockChain{
		tip: tip,
		db:  db,
	}, nil
}

func (bc *BlockChain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		lastHash = b.Get([]byte("l"))
		return nil
	})

	if err != nil {
		log.Println("Error starting a read only transaction")
		os.Exit(0)
	}

	newBlock := NewBlock(data, lastHash)

	bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		b.Put(newBlock.Hash, newBlock.Serialize())
		b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash
		return nil
	})
}
