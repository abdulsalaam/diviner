# Logarithmic Market Scoring Rule (LMSR) Prediction Market on Hyperledger Fabric

The project is a prototype implements **Prediction Market** with **LMSR** on IBM **Hyperledger Fabric**. Fabric is a block chain open source like Bitcoin or Ethereum.

## Resources
1. [Logarithmic Market Scoring Rules for Modular Combinatorial Information Aggregation by Robin Hanson](http://mason.gmu.edu/~rhanson/mktscore.pdf)
2. [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/release/)
3. [Gnosis](https://gnosis.pm/)
4. [Prediction Market (en)](https://en.wikipedia.org/wiki/Prediction_market)
5. [Prediction Market (ch)](https://zh.wikipedia.org/wiki/%E9%A2%84%E6%B5%8B%E5%B8%82%E5%9C%BA)

## Fixtures

### preparation
1. `cryptogen generate --config=./crypto-config.yaml`
2. `export FABRIC_CFG_PATH=$PWD`
3. `mkdir channel-artifacts`
4. `configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block`
5. `configtxgen -profile DivinerChannel -outputCreateChannelTx ./channel-artifacts/diviner_channel.tx -channelID divinerchannel`
6. `configtxgen -profile DivinerChannel -outputAnchorPeersUpdate ./channel-artifacts/DivinerMSPanchors.tx -channelID divinerchannel -asOrg DivinerMSP`

### start
1. `docker-compose -f docker-compose-cli.yaml up -d`

### stop and remove
1. `docker-compose -f docker-compose-cli.yaml down`

### login cli
1. `docker exec -it cli bash`

### cli
#### join channel
1. `peer channel create -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/diviner_channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem`
2. `peer channel join -b divinerchannel.block`

3. `CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/users/Admin@diviner.info/msp CORE_PEER_ADDRESS=peer1.diviner.info:7051 CORE_PEER_LOCALMSPID="DivinerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/peers/peer1.diviner.info/tls/ca.crt peer channel join -b divinerchannel.block`

#### update anchor
1. `peer channel update -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/DivinerMSPanchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem`

#### chaincode
0. `go get -u github.com/golang/dep/cmd/dep`
0. `dep ensure`
1. `peer chaincode install -n member -v 1.0 -p diviner/chaincode/member`
2. `peer chaincode instantiate -o orderer.diviner.info:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem -C divinerchannel -n member -v 1.0 -c '{"Args":[]}' -P "OR ('DivinerMSP.member')"`

## TOFix
1. `Error: Error endorsing chaincode: rpc error: code = Unknown desc = error starting container: API error (400): {"message":"oci runtime error: container_linux.go:265: starting container process caused \"exec: \\\"chaincode\\\": executable file not found in $PATH\"\n"}` after upgrading fabric 1.1.0-preview docker images

## TODO
1. Because can not handle concurrent transactions on a market, it needs an transaction queue for each market
2. Member management
3. Mobile App
