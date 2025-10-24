package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
)

func (cli *CLI) getBalance(address string) {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicln("Error opening db file")
	}

	var genesisBlockHash []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Println("No bucket found with name: ", bucket)
			return err
		}
		genesisBlockHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panicln("Error starting a read only transaction")
	}

	bc := &Blockchain{
		tip: genesisBlockHash,
		db:  db,
	}
	cli.bc = bc

	unspentoutputs := cli.bc.FindUTXOs(address)
	balance := 0

	for _, output := range unspentoutputs {
		balance += output.Value
	}

	fmt.Printf("Balance of %s: %v\n", address, balance)

	err = bc.db.Close()
	if err != nil {
		log.Panicln("Error closing the db while fetching the balance")
	}
}

func (cli *CLI) printChain() {
	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicln("Error opening db file")
	}

	var genesisBlockHash []byte

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Println("No bucket found with name: ", bucket)
			return err
		}
		genesisBlockHash = b.Get([]byte("l"))
		return nil
	})
	if err != nil {
		log.Panicln("Error starting a read only transaction")
	}

	// fmt.Printf("Genesis hash: %x\n", genesisBlockHash)

	bc := &Blockchain{
		tip: genesisBlockHash,
		db:  db,
	}
	cli.bc = bc

	iterator := cli.bc.Iterator()

	for {
		block := iterator.Next()
		fmt.Println()
		fmt.Printf("Previous Hash: %x\n", block.PreviousBlockHash)
		// fmt.Printf("Data: %v\n", block.Transactions)
		fmt.Printf("Current Hash: %x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("Proof of work: %v\n", pow.Validate())
		fmt.Println()

		if len(block.PreviousBlockHash) == 0 {
			break
		}
	}
}

func (cli *CLI) createBlockchain(address string) {
	bc := CreateBlockChain(address)
	err := bc.db.Close()
	if err != nil {
		log.Panicln("Error closing the db")
	}
	fmt.Println("Done")
}

func (cli *CLI) send(from, to string, amount int) {
	bc, err := NewBlockChain(from)
	if err != nil {
		log.Panicln("Error initializing blockchain")
	}
	defer func() {
		if err := bc.db.Close(); err != nil {
			log.Panicln("Error closing the database")
		}
	}()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success")
}

func (cli *CLI) createWallet() {
	wallets := NewWallets()
	address := wallets.CreateWallet()
	fmt.Println("New address: ", address)
}
