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

## TODO
1. Because can not handle concurrent transactions on a market, it needs an transaction queue for each market
2. Member management
3. Mobile App
