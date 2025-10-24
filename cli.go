package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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

func (cli *CLI) Run() {
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createBlockchainCmd := flag.NewFlagSet("createblockchain", flag.ExitOnError)
	getBalCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
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
	case "createblockchain":
		if err := createBlockchainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the createblockchain command")
		}
	case "getbalance":
		if err := getBalCmd.Parse(os.Args[2:]); err != nil {
			log.Panicln("Error parsing the getbalance command")
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
