#!/usr/bin/env bash
set -ex


echo "Running chainbridge..."
./build/chainbridge-celo --config ./e2e/config-celo-int-tst.json --testkey alice --fresh --leveldb ./lvldb
# Otherwise CI will run tests before ganache has started
sleep 15