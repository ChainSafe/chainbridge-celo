# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only
FROM chainsafe/celo-chaindata:latest AS dump

WORKDIR /src

FROM us.gcr.io/celo-org/geth:1.2.2

WORKDIR /root

COPY --from=dump /src/ .
COPY entrypoint.sh .
COPY ./password.txt .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["/root/entrypoint.sh"]
