#!/usr/bin/env bash
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

set -eux

cd ${SCRIPT_DIR}/celo/elections

if ! $(ls | grep -q node_modules) ;
then
    npm install
fi

npm start

