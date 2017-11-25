package main

import (
	"diviner/chaincode"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(chaincode.NewEventChaincode())

	if err != nil {
		fmt.Printf("creating event chaincode failed: %v", err)
	}
}
