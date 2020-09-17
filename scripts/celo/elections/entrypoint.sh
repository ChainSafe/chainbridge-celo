#!/usr/bin/env sh
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

# Exit on failure
set -ex

geth init /root/genesis.json --datadir /root/.celo1/
rm -rf /root/.celo1/keystore
cp -r /root/keystore /root/.celo1/

geth init /root/genesis.json --datadir /root/.celo2/
rm -rf /root/.celo2/keystore
cp -r /root/keystore /root/.celo2/
rm /root/.celo2/keystore/UTC--2020-03-20T21-14-07.508717834Z--f4314cb9046bece6aa54bb9533155434d0c76909

exec geth \
  --datadir /root/.celo1 \
  --unlock "0xf4314cb9046bece6aa54bb9533155434d0c76909","0xff93B45308FD417dF303D6515aB04D9e89a750Ca","0x8e0a907331554AF72563Bd8D43051C2E64Be5d35","0x24962717f8fA5BA3b931bACaF9ac03924EB475a0","0x148FfB2074A9e59eD58142822b3eB3fcBffb0cd7","0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485" \
  --password /root/password.txt \
  --port 30303 \
  --ws \
  --wsport 8545 \
  --wsorigins="*" \
  --wsaddr 0.0.0.0 \
  --rpc \
  --rpcport 8546 \
  --rpccorsdomain="*" \
  --networkid 160693804425899027 \
  --rpcaddr 0.0.0.0 \
  --allow-insecure-unlock \
  --mine \
  --nodekeyhex 98ab333347a12cef869f92b3de44085f9e44891e513bcf1d76a99eecbcdd5e17 \
  --bootnodes enode://20172c242c53d46ef9f07c2cd046e86719b7f99107eea995dc5026bc867448c455536a123d430353593125a2adb3d1a0ebb83b7f85228df3667c0bd941a107d0@127.0.0.1:30304 \
  --syncmode full &

exec geth \
  --datadir /root/.celo2 \
  --unlock "0xff93B45308FD417dF303D6515aB04D9e89a750Ca","0x8e0a907331554AF72563Bd8D43051C2E64Be5d35","0x24962717f8fA5BA3b931bACaF9ac03924EB475a0","0x148FfB2074A9e59eD58142822b3eB3fcBffb0cd7","0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485" \
  --password /root/password.txt \
  --port 30304 \
  --ws \
  --wsport 8645 \
  --wsorigins="*" \
  --wsaddr 0.0.0.0 \
  --rpc \
  --rpcport 8646 \
  --rpccorsdomain="*" \
  --networkid 160693804425899027 \
  --rpcaddr 0.0.0.0 \
  --allow-insecure-unlock \
  --mine \
  --nodekeyhex 98ab333347a12cef869f92b3de44085f9e44891e513bcf1d76a99eecbcdd5e18 \
  --bootnodes "enode://47eccabfa6fa81f7ae7d8c094485fd521789e22e529916a2409970b6e99b1d1b6c590ad71bd08f182071313b2e03d319f4f6b625b28b6ab7768bb2de3f2a6eed@127.0.0.1:30303" \
  --syncmode full
