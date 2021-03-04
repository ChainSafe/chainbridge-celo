#!/usr/bin/env bash
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

CMD=chainbridge-celo

BRIDGE_ADDRESS="0x62877dDCd49aD22f5eDfc6ac108e9a4b5D2bD88B"
ERC20_ADDRESS="0x21605f71845f372A9ed84253d2D024B7B10999f4"
ERC20_HANDLER="0x3167776db165D8eA0f51790CA2bbf44Db5105ADF"
RESOURCE_ID="000000000000000000000021605f71845f372A9ed84253d2D024B7B10999f400"

GAS_LIMIT=6721975
GAS_PRICE=20000000000


ERC721_HANDLER="0x3167776db165D8eA0f51790CA2bbf44Db5105ADF"
ERC721_RESOURCE_ID="0000000000000000000000d7E33e1bbf65dC001A0Eb1552613106CD7e40C3100"
ERC721_CONTRACT="0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"

GENERIC_HANDLER="0x106C24dc2D480b5559C9E0e97bAaDf0750d9F0B8"
GENERIC_RESOURCE_ID="0000000000000000000000106C24dc2D480b5559C9E0e97bAaDf0750d9F0B800"


set -eux

#deploy
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --all --erc20Symbol "TKN" --erc20Name "token  token"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --erc721
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --erc20 --erc20Symbol "TKN" --erc20Name "token  token"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE deploy --bridge


#erc20
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 mint --amount 100 --erc20Address $ERC20_ADDRESS
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 add-minter --erc20Address $ERC20_ADDRESS --minter "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 allowance --erc20Address $ERC20_ADDRESS --spender "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E" --owner "0x2f709398808af36ADBA86ACC617FeB7F5B7B1931"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 approve --erc20Address $ERC20_ADDRESS --recipient "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E" --amount "1.11"  --decimals 2
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc20 balance --erc20Address $ERC20_ADDRESS --address "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
#
#
#admin
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin add-admin --bridge $BRIDGE_ADDRESS --admin "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin add-relayer --bridge $BRIDGE_ADDRESS --relayer "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin is-relayer --bridge $BRIDGE_ADDRESS --relayer "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin pause --bridge $BRIDGE_ADDRESS
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin remove-admin --bridge $BRIDGE_ADDRESS --admin "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin remove-relayer --bridge $BRIDGE_ADDRESS --relayer "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin set-fee --bridge $BRIDGE_ADDRESS --fee 123 --decimals 10
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin set-threshold --bridge $BRIDGE_ADDRESS --threshold 2
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE admin unpause --bridge $BRIDGE_ADDRESS
#
#brdige
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge query-resource --handler $ERC20_HANDLER --resourceId $RESOURCE_ID
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge set-burn --bridge $BRIDGE_ADDRESS  --handler $ERC20_HANDLER --tokenContract $ERC20_ADDRESS
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge register-resource --bridge $BRIDGE_ADDRESS  --handler $ERC20_HANDLER --resourceId $RESOURCE_ID --targetContract $ERC20_ADDRESS
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE bridge register-generic-resource --bridge "0x39863e3eDB5255dB93bBf8E76c12578357dBe6c7"  --handler "0x106C24dc2D480b5559C9E0e97bAaDf0750d9F0B8" --resourceId "0000000000000000000000106C24dc2D480b5559C9E0e97bAaDf0750d9F0B800" --targetContract $ERC20_ADDRESS --hash true --execute "transfer(address,uint256)"
#
#
#
#erc721
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc721 --erc721Address $ERC721_CONTRACT mint --id 0x1
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc721 --erc721Address $ERC721_CONTRACT add-minter --minter "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc721 --erc721Address $ERC721_CONTRACT approve --id 0x1 --recipient "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E"
$CMD cli --gasLimit $GAS_LIMIT --gasPrice $GAS_PRICE erc721 deposit --id 0x1 --bridge $BRIDGE_ADDRESS --dest 5 --recipient "0x3f709398808af36ADBA86ACC617FeB7F5B7B193E" --resourceId $ERC721_RESOURCE_ID
