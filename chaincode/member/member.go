package main

import (
	"diviner/chaincode"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(chaincode.NewMemberChaincode())

	if err != nil {
		fmt.Printf("creating member chaincode failed: %v", err)
	}
}
