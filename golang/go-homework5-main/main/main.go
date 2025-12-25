package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

func main() {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000000")
	fmt.Println(addr.Hex())

}
