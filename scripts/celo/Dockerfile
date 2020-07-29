# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only
FROM us.gcr.io/celo-testnet/celo-node:alfajores

WORKDIR /root

COPY ./genesis.json .
COPY ./keystore /root/keystore
COPY entrypoint.sh .
COPY ./password.txt .

RUN chmod +x entrypoint.sh

ENTRYPOINT ["/root/entrypoint.sh"]