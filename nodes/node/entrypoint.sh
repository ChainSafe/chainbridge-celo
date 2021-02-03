#!/usr/bin/env sh
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

# Exit on failure
set -ex

#geth init /root/genesis.json
#rm -f /root/.celo/keystore/*

# If accounts are not set, set all accounts.
if [ -z $ACCOUNTS ]
then
  ACCOUNTS="0xf4314cb9046bece6aa54bb9533155434d0c76909,0xff93B45308FD417dF303D6515aB04D9e89a750Ca,0x8e0a907331554AF72563Bd8D43051C2E64Be5d35,0x24962717f8fA5BA3b931bACaF9ac03924EB475a0,0x148FfB2074A9e59eD58142822b3eB3fcBffb0cd7,0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"
fi
# Copy requested accounts to celo keystore.

# Identify the docker container external IP.
IP=$(ip -4 -o address | \
  grep -Eo -m 1 'eth0\s+inet\s+[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}' | \
  grep -Eo '[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}')

if [ ! -z $BOOTNODE ]
then
  BOOTNODE="--bootnodes ${BOOTNODE}"
else
  BOOTNODE="--nodiscover"
fi
if [ ! -z $NODEKEY ]; then NODEKEY="--nodekeyhex ${NODEKEY}"; fi
if [ -z $NETWORKID ]; then NETWORKID="5"; fi

if [ -z $DATADIR ]; then DATADIR="/root/celo-dump1"; fi

exec geth \
  --unlock ${ACCOUNTS} \
  --password password.txt \
  --datadir=${DATADIR} \
  --ws \
  --wsport 8545 \
  --wsorigins="*" \
  --wsaddr 0.0.0.0 \
  --rpc \
  --rpcport 8546 \
  --rpccorsdomain="*" \
  --rpcvhosts="*" \
  --nat=extip:${IP} \
  --networkid ${NETWORKID} \
  --rpcaddr 0.0.0.0 \
  --allow-insecure-unlock \
  --mine ${BOOTNODE} ${NODEKEY}

