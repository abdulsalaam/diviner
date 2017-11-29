#!/bin/bash

peer channel create -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/diviner_channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem
peer channel join -b divinerchannel.block
CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/users/Admin@diviner.info/msp CORE_PEER_ADDRESS=peer1.diviner.info:7051 CORE_PEER_LOCALMSPID="DivinerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/peers/peer1.diviner.info/tls/ca.crt peer channel join -b divinerchannel.block
peer channel update -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/DivinerMSPanchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem

# chaincodes
peer chaincode install -n lmsr -v 1.0 -p diviner/chaincode
peer chaincode instantiate -o orderer.diviner.info:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem -C divinerchannel -n lmsr -v 1.0 -c '{"Args":[]}' -P "OR ('DivinerMSP.member')"
