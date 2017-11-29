# Logarithmic Market Scoring Rule (LMSR) Prediction Market on Hyperledger Fabric

The project is a prototype implements **Prediction Market** with **LMSR** on IBM **Hyperledger Fabric**. Fabric is a block chain open source like Bitcoin or Ethereum.

## Resources
1. [Logarithmic Market Scoring Rules for Modular Combinatorial Information Aggregation by Robin Hanson](http://mason.gmu.edu/~rhanson/mktscore.pdf)
2. [Hyperledger Fabric](https://hyperledger-fabric.readthedocs.io/en/release/)
3. [Gnosis](https://gnosis.pm/)
4. [Prediction Market (en)](https://en.wikipedia.org/wiki/Prediction_market)
5. [Prediction Market (ch)](https://zh.wikipedia.org/wiki/%E9%A2%84%E6%B5%8B%E5%B8%82%E5%9C%BA)

## How to run it
1. clone this project: `git clone git@github.com:kigichang/diviner.git`
2. go to folder **diviner/fixtures**
3. start fabric environment: `./start.sh`
4. watch **cli** (run `docker logs cli`) logs until chaincode is instantiated (you will see `[main] main -> INFO 009 Exiting.....`)
5. open a new terminal, go to **diviner/service**, and start middle service: `go run diviner_service.go`
6. open a new terminal, go to **diviner/app**, and play. :)

### How to play
1. create member: `go run app.go member create --ski [private key]`
  * `go run app.go member create --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
  * `go run app.go member create --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

2. create event: `go run app.go event create --title title --outcome outcome1 --outcome outcome2 --ski [private key]`
  * `go run app.go event create --title gogo --outcome yes --outcome no --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`

    ps: remember the event id on screen like: `event:<id:"4f2107e3-6aa8-45b0-8419-fe8c492756d5" .... >` (the event id is *4f2107e3-6aa8-45b0-8419-fe8c492756d5*)

3. create market: `go run app.go market create --event [event_id got from  step2] --number 100 --ski [private key]`
  * `go run app.go market create --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1 --event 4f2107e3-6aa8-45b0-8419-fe8c492756d5 --number 100`

    ps: *a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1* is the event id got from step2

    ps: remember the market and share id on screen like: ` market:<id:"4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed" ... event:"4f2107e3-6aa8-45b0-8419-fe8c492756d5" liquidity:100 fund:69.31471805599453 cost:69.31471805599453 shares:<key:"4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@0" value:0 > shares:<key:"4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@1" value:0 > ... >` (the market id is *4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed*, **yes** share id is *4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@0*, and **no** share id is *4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@1*)

4. buy or sell a ahare: `go run app.go tx [buy or sell] [share id got from step3] [volume] --ski [private key]`

  * `go run app.go tx buy 4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@0 100 --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

    ps. *4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@0* is the **yes** share id got from step3

5. query members and will see the balance and assets changed.
  * `go run app.go member query --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`

    the balance is *99930.68528194401*

  * `go run app.go member query --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

    the balance is *99937.98854930417* and have asset *4f2107e3-6aa8-45b0-8419-fe8c492756d5#58c77b4c-f686-4170-a78a-0d96308496ed#4f2107e3-6aa8-45b0-8419-fe8c492756d5@0#aSq9DsNNvGhYxYyqA9wd2eduEAZ5AXWgJTbTJAb2Wjq7bj9GGAWwADm9W2UknFWmZhqxd2G21L9WTXyxSHfeR8z1moeVjHNjEFmRqHFRBPC8ckqwSfsKJZc814kR* with volume *100*

### Clean all data
1. go to **diviner/fixtures** and run `./teardown.sh` to clean all container and dev images.

### Other commands
* Query Member: `go run app.go member query --ski [private key]`
  1. `go run app.go member query --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`
  2. `go run app.go member query --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`


* Query Event: `go run app.go event query --ski [private key]`
  1. `go run app.go event query --id [event id] --ski a7145e9a7b7bea5907bb022333beaac24bc4095d17f417f262b543de2c54bed1`


* Query Market: `go run app.go market query --id [market id]`
  1. `go run app.go market query --id a124ddec-a673-437c-be5c-0c54bdf58366#3981d198-ca8f-4f33-a0aa-19a6466ce984 --ski f306ad8811b1d1649bf96d3faeffd7d0c3a21a1fc855481adc7b9be52c596ed6`

## Fixtures detail

Because all files are in folder, you may not run all steps in **preparation**

### preparation
1. `cryptogen generate --config=./crypto-config.yaml`
2. `export FABRIC_CFG_PATH=$PWD`
3. `mkdir channel-artifacts`
4. `configtxgen -profile TwoOrgsOrdererGenesis -outputBlock ./channel-artifacts/genesis.block`
5. `configtxgen -profile DivinerChannel -outputCreateChannelTx ./channel-artifacts/diviner_channel.tx -channelID divinerchannel`
6. `configtxgen -profile DivinerChannel -outputAnchorPeersUpdate ./channel-artifacts/DivinerMSPanchors.tx -channelID divinerchannel -asOrg DivinerMSP`


The cli container will run **scripts/script.sh** to initialize all environment when starting. The detail is following:

### join channel
1. `peer channel create -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/diviner_channel.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem`
2. `peer channel join -b divinerchannel.block`

3. `CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/users/Admin@diviner.info/msp CORE_PEER_ADDRESS=peer1.diviner.info:7051 CORE_PEER_LOCALMSPID="DivinerMSP" CORE_PEER_TLS_ROOTCERT_FILE=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/diviner.info/peers/peer1.diviner.info/tls/ca.crt peer channel join -b divinerchannel.block`

#### update anchor
1. `peer channel update -o orderer.diviner.info:7050 -c divinerchannel -f ./channel-artifacts/DivinerMSPanchors.tx --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem`

#### chaincode
1. `peer chaincode install -n lmsr -v 1.0 -p diviner/chaincode`
2. `peer chaincode instantiate -o orderer.diviner.info:7050 --tls $CORE_PEER_TLS_ENABLED --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/diviner.info/orderers/orderer.diviner.info/msp/tlscacerts/tlsca.diviner.info-cert.pem -C divinerchannel -n lmsr -v 1.0 -c '{"Args":[]}' -P "OR ('DivinerMSP.member')"`

## TODO
1. Because can not handle concurrent transactions on a market, it needs an transaction queue for each market
2. Management
3. Mobile App

## Docker Others
* start: `docker-compose -f docker-compose-cli.yaml up -d`
* stop: `docker-compose -f docker-compose-cli.yaml down`
* login cli: `docker exec -it cli bash`
* rm containers: `docker rm $(docker ps -aq)`
* rm dev images: `docker rmi $(docker images --filter=reference='dev*' -q)`
* rm all images: `dokcer rmi $(docker images -q)`
