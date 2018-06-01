#!/bin/bash

rm -rf ./artifacts/*
rm -rf ./crypto-config/*

cryptogen generate --config=./crypto-config.yaml

FABRIC_CFG_PATH=$PWD configtxgen -profile ChainHero -outputBlock ./artifacts/orderer.genesis.block

FABRIC_CFG_PATH=$PWD configtxgen -profile ChainHero -outputCreateChannelTx ./artifacts/chainhero.channel.tx -channelID chainhero

FABRIC_CFG_PATH=$PWD configtxgen -profile ChainHero -outputAnchorPeersUpdate ./artifacts/org1.chainhero.anchors.tx -channelID chainhero -asOrg ChainHeroOrganization1

