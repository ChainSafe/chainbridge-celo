#!/usr/bin/env bash
set -ex


echo "Running chainbridge..."
./build/chainbridge-celo run --config ./e2e/cfg/config-celo-int-tst.json --testkey alice --fresh --leveldb ./lvldb &
# Otherwise CI will run tests before chain has started
sleep 15
