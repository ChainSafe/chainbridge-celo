#!/usr/bin/env sh
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

# Exit on failure
set -ex

geth init /root/genesis.json

rm -rf /root/.celo/keystore
cp -r /root/keystore /root/.celo/

exec geth \
  --unlock "0xf4314cb9046bece6aa54bb9533155434d0c76909","0xff93B45308FD417dF303D6515aB04D9e89a750Ca","0x8e0a907331554AF72563Bd8D43051C2E64Be5d35","0x24962717f8fA5BA3b931bACaF9ac03924EB475a0","0x148FfB2074A9e59eD58142822b3eB3fcBffb0cd7","0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485" \
  --password /root/password.txt \
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
  --bootnodes enode://20172c242c53d46ef9f07c2cd046e86719b7f99107eea995dc5026bc867448c455536a123d430353593125a2adb3d1a0ebb83b7f85228df3667c0bd941a107d0@10.11.1.2:30303 \
  --nat extip:10.11.1.1 \
  --syncmode full
