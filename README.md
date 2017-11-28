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
1. `peer chaincode install -n lmsr -v 1.0 -p diviner/chaincode`
2. `peer chaincode instantiate -o orderer.diviner.info:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem -C divinerchannel -n lmsr -v 1.0 -c '{"Args":[]}' -P "OR ('DivinerMSP.member')"`

## App
### member
1. `go run app.go member create --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
2. `go run app.go member create --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`
3. `go run app.go member query --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
4. `go run app.go member query --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

### event
1. `go run app.go event create --title gogo --outcome yes --outcome no --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
2. `go run app.go event query --id a124ddec-a673-437c-be5c-0c54bdf58366 --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`

## market
1. `go run app.go market create --event a124ddec-a673-437c-be5c-0c54bdf58366 --number 10000 --fund --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
2. `go run app.go market query --id a124ddec-a673-437c-be5c-0c54bdf58366#3981d198-ca8f-4f33-a0aa-19a6466ce984 --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

## TODO
1. Because can not handle concurrent transactions on a market, it needs an transaction queue for each market
2. Member management
3. Mobile App

## Docker Others
* start: `docker-compose -f docker-compose-cli.yaml up -d`
* stop: `docker-compose -f docker-compose-cli.yaml down`
* login cli: `docker exec -it cli bash`
* rm containers: `docker rm $(docker ps -aq)`
* rm dev images: `docker rmi $(docker images --filter=reference='dev*' -q)`
* rm all images: `dokcer rmi $(docker images -q)`
