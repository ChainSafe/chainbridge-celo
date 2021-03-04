# Admin Command

- [`is-relayer`](#is-relayer)
- [`add-relayer`](#add-relayer)
- [`remove-relayer`](#remove-relayer)
- [`set-threshold`](#set-threshold)
- [`pause`](#pause)
- [`unpause`](#unpause)
- [`set-fee`](#set-fee)
- [`withdraw`](#withdraw)
- [`add-admin`](#add-admin)
- [`remove-admin`](#remove-admin)

## `is-relayer`
Check if an address is registered as a relayer.

```
--relayer <value>   Address to check
--bridge <address>  Bridge contract address
```

## `add-relayer`
Adds a new relayer.

```
--relayer <address>  Address of relayer
--bridge <address>   Bridge contract address
```

## `remove-relayer`
Removes a relayer.

```
--relayer <address>  Address of relayer
--bridge <address>   Bridge contract address
```

## `set-threshold`
Sets a new relayer vote threshold.

```
--bridge <address>   Bridge contract address
--threshold <value>  New relayer threshold
```

## `pause`
Pauses deposits and proposals.

```
--bridge <address>  Bridge contract address 
```

## `unpause`
Unpause deposits and proposals.

```
--bridge <address>  Bridge contract address 
```

## `set-fee`
Set a new fee.

```
--bridge <address>   Bridge contract address
--fee <string>       New fee in ethers (can be float)
--decimals           The number of decimal places for the erc20 token
```

## `withdraw`
Withdraw tokens from a handler contract.

```
--bridge <address>         Bridge contract address
--handler <address>        Handler contract address
--tokenContract <address>  ERC20 or ERC721 token contract address
--recipient <address>      Address to withdraw to
--amount <string>          Token amount to withdraw. Should be set or id or amount if both set error will occur
--id <uint64>              ID of token to withdraw. Should be set or id or amount if both set error will occur
--decimals                 The number of decimal places for the erc20 token. Provide only when amount is specified

```

## `add-admin`
Adds an admin

```
--admin <address>   Address of admin
--bridge <address>  Bridge contract address
```

## `remove-admin`
Removes an admin

```
--admin <address>   Address of admin
--bridge <address>  Bridge contract address
```

