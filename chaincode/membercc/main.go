package main

import (
	"diviner/chaincode/member"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

func main() {
	err := shim.Start(member.NewMemberChaincode())
	if err != nil {
		fmt.Printf("creating member chaincode failed: %v\n", err)
	}
}
