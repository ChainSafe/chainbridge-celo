#!/usr/bin/env sh
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

# Exit on failure
set -ex

geth init /root/genesis.json
rm -f /root/.celo/keystore/*

ACCOUNTS_TRIMMED=${ACCOUNTS//0x/}
ACCOUNTS_PATTERN=${ACCOUNTS_TRIMMED//,/|}
find /root/keystore | grep -iE ${ACCOUNTS_PATTERN} | xargs -i cp {} /root/.celo/keystore/

IP=$(ip -4 -o address | \
  grep -Eo -m 1 'eth0\s+inet\s+[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}' | \
  grep -Eo '[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}')

exec geth \
  --unlock ${ACCOUNTS} \
  --password /root/password.txt \
  --ws \
  --wsport 8545 \
  --wsorigins="*" \
  --wsaddr 0.0.0.0 \
  --rpc \
  --rpcport 8546 \
  --rpccorsdomain="*" \
  --rpcvhosts="*" \
  --nat=extip:${IP} \
  --v5disc \
  --networkid 160693804425899027 \
  --rpcaddr 0.0.0.0 \
  --allow-insecure-unlock \
  --mine \
  --syncmode full \
  --bootnodes ${BOOTNODE} \
  --nodekeyhex ${NODEKEY}

