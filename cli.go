package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

type CLI struct {
	bc *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

// func (cli *CLI) addBlock(data string) {
// 	cli.bc.MineBlock(data)
// 	fmt.Println("Success")
// }

func (cli *CLI) printChain() {
	db, err :=	bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panicln("Error opening db file")
	}
	
	var genesisBlockHash []byte

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			log.Panicln("No bucket found with name: ", bucket)
		}
		genesisBlockHash = b.Get([]byte("l"))
		return nil
	})

	// fmt.Printf("Genesis hash: %x\n", genesisBlockHash)

	bc := &Blockchain{
		tip: genesisBlockHash,
		db: db,
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
	bc.db.Close()
	fmt.Println("Done")
}

func (cli *CLI) Run() {
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockChainCmd := flag.NewFlagSet("createBlockChain", flag.ExitOnError)

	createBlockChainAddress := createBlockChainCmd.String("address", "", "The address to send the reward for mining the genesis block")

	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "printChain":
		printChainCmd.Parse(os.Args[2:])
	case "createBlockChain":
		createBlockChainCmd.Parse(os.Args[2:])
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockChainCmd.Parsed() {
		if *createBlockChainAddress == "" {
			createBlockChainCmd.Usage()
			os.Exit(1)
		}

		cli.createBlockchain(*createBlockChainAddress)
	}
}
