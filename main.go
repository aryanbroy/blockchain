package main

import (
// "os"
)

func main() {
	// bc, err := NewBlockChain()
	// if err != nil {
	// 	os.Exit(0)
	// }

	// defer bc.db.Close()

	// cli := CLI{bc: bc}
	cli := CLI{}
	cli.Run()
}
