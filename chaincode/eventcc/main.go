package main

import (
	"diviner/chaincode/event"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {

	err := shim.Start(event.NewEventChaincode())
	if err != nil {
		fmt.Printf("creating event chaincode failed: %v\n", err)
	}
}
