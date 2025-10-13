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
	fmt.Println("  createblockchain -address ADDRESS - Create a blockchain and send genesis block reward to ADDRESS")
	fmt.Println("  createwallet - Generates a new key-pair and saves it into the wallet file")
	fmt.Println("  getbalance -address ADDRESS - Get balance of ADDRESS")
	fmt.Println("  listaddresses - Lists all addresses from the wallet file")
	fmt.Println("  printchain - Print all the blocks of the blockchain")
	fmt.Println("  send -from FROM -to TO -amount AMOUNT - Send AMOUNT of coins from FROM address to TO")
}

// func (cli *CLI) addBlock(data string) {
// 	cli.bc.MineBlock(data)
// 	fmt.Println("Success")
// }
//

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

func (cli *CLI) Run() {
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createBlockchain", flag.ExitOnError)
	getBalCmd := flag.NewFlagSet("getBalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)
	createWalletCmd := flag.NewFlagSet("createwallet", flag.ExitOnError)

	createBlockchainAddress := createBlockchainCmd.String("address", "", "The address to send the reward for mining the genesis block")
	getBalAddress := getBalCmd.String("address", "", "Address to fetch balance for")
	from := sendCmd.String("from", "", "Sender of the transaction")
	to := sendCmd.String("to", "", "Recipent of the transaction")
	amount := sendCmd.Int("amount", 0, "Amount to be sent")

	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "createwallet":
		if err := createWalletCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the createwallet command")
		}
	case "printChain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the printChain command")
		}
	case "createBlockchain":
		if err := createBlockchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the createblockchain command")
		}
	case "getBalance":
		if err := getBalCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the getBalance command")
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the send command")
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if createWalletCmd.Parsed() {
		cli.createWallet()
	}

	if sendCmd.Parsed() {
		cli.send(*from, *to, *amount)
	}

	if getBalCmd.Parsed() {
		cli.getBalance(*getBalAddress)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}

	if createBlockchainCmd.Parsed() {
		if *createBlockchainAddress == "" {
			createBlockchainCmd.Usage()
			os.Exit(1)
		}

		cli.createBlockchain(*createBlockchainAddress)
	}
}

func (cli *CLI) createWallet() {
	_ = NewWallets()

	// wallet := NewWallet()
	// address := wallet.GenerateAddress()
	// fmt.Printf("Wallet address: %x\n", address)
}
