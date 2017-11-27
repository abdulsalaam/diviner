package main

import (
	"diviner/chaincode"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	logger := shim.NewLogger("event")
	logger.SetLevel(shim.LogDebug)
	logger.Info("start event chaincode")

	err := shim.Start(chaincode.NewEventChaincode())

	if err != nil {
		fmt.Printf("creating event chaincode failed: %v", err)
	}
}
