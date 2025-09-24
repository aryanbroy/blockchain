package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type Blockchain struct {
	tip []byte
	db  *bolt.DB
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

const dbFile = "blockchain.db"
const bucket = "nigaMania"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

func dbExists(db string) bool {
	if _, err := os.Stat(db); os.IsNotExist(err) {
		return false
	}

	return true
}

func NewBlockChain(address string) (*Blockchain, error) {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Println("Error initializing database")
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			coinbaseTx := NewCoinbaseTx(address, genesisCoinbaseData)
			genesis := NewGenesisBlock(coinbaseTx)

			b, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				log.Println("Error creating a bucket")
				return err
			}

			if err := b.Put(genesis.Hash, genesis.Serialize()); err != nil {
				log.Println("Error updating data in db")
				return err
			}

			if err := b.Put([]byte("l"), genesis.Hash); err != nil {
				log.Println("Error updating data in db")
				return err
			}
			tip = genesis.Hash

		} else {
			tip = b.Get([]byte("l"))
		}
		return nil
	})
	if err != nil {
		log.Println("Error updating the database.", err)
		return nil, err
	}

	return &Blockchain{
		tip: tip,
		db:  db,
	}, nil
}

func CreateBlockChain(address string) *Blockchain {
	if dbExists(address) {
		log.Println("Blockchain already exists")
		os.Exit(1)
	}

	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Println("Error opening the db file")
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucket([]byte(bucket))
		if err != nil {
			if err == bolt.ErrBucketExists {
				err := fmt.Errorf("bucket already exists")
				return err
			}
			log.Println("Error creating a bucket")

			return err
		}

		log.Println("Generating a new genesis block...")
		fmt.Println()

		coinbaseTx := NewCoinbaseTx(address, genesisCoinbaseData)
		genesisBlock := NewGenesisBlock(coinbaseTx)

		err = b.Put(genesisBlock.Hash, genesisBlock.Serialize())
		if err != nil {
			log.Println("Error updating db values")
			return err
		}

		err = b.Put([]byte("l"), genesisBlock.Hash)
		if err != nil {
			log.Println("Error updating last blockchain value")
			return err
		}

		tip = genesisBlock.Hash
		return nil
	})

	if err != nil {
		log.Println("Error creating a new blockchain")
		log.Panic(err)
	}

	return &Blockchain{
		tip: tip,
		db:  db,
	}

}

func (bc *Blockchain) MineBlock(transactions []*Transaction) {
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

	newBlock := NewBlock(transactions, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			log.Println("error updating the db")
			return err
		}
		if err := b.Put([]byte("l"), newBlock.Hash); err != nil {
			log.Println("error updating the db")
			return err
		}
		bc.tip = newBlock.Hash
		return nil
	})

	if err != nil {
		log.Println("error mining the block")
		os.Exit(1)
	}
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{
		currentHash: bc.tip,
		db:          bc.db,
	}
}

func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		encodedBlock := b.Get(i.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})

	if err != nil {
		log.Println("Error starting a transaction!")
		os.Exit(0)
	}

	i.currentHash = block.PreviousBlockHash

	return block
}
