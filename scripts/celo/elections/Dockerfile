# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only
FROM node:12.18.4

WORKDIR /root

COPY contracts ./contracts/
COPY index.js .
COPY package.json .
COPY package-lock.json .
COPY entrypoint.sh .
RUN npm install

RUN chmod +x entrypoint.sh

ENTRYPOINT ["/root/entrypoint.sh"]
