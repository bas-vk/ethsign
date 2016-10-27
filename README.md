## Ethereum signer
Ethsigner is a utility that uses local keystore and can sign transactions using an interactive console.


### Usage
Ethsigner provides the following features:

`eth.accounts` lists accounts in the keystore

`personal.newAccount([password])` create new account, if password is not given you will be prompted

`personal.unlockAccount(<address>, [password], [duration])` unlock account, if password is not given you will be prompted.

`eth.signTransaction({transaction object})` sign transaction

### Example
1. Start ethsign
```
ethsign --keystore </path/to/keystore>
```

2. Create account
```
> personal.newAccount()
Passphrase: 
Repeat passphrase: 
"0x58c5a268e50edbe3040969b6e742fd5af4c9e412"
```

2. List accounts
```
> eth.accounts
["0x58c5a268e50edbe3040969b6e742fd5af4c9e412"]
```

3. Unlock account for 1 minute
```
> personal.unlockAccount("0x58c5a268e50edbe3040969b6e742fd5af4c9e412", null, 60)
Unlock account 0x58c5a268e50edbe3040969b6e742fd5af4c9e412
Passphrase: 
true
```

3. Sign value transaction
```
> eth.signTransaction({from: "0x58c5a268e50edbe3040969b6e742fd5af4c9e412", to: "0x0f13c906944b155ef6ed569b1ce72a385b090fdc", value: 1, nonce: 12, gas: 21000, gasPrice: web3.toWei(20, "shannon")})
"0xf8640c8504a817c800825208940f13c906944b155ef6ed569b1ce72a385b090fdc01801ca0dcbb20010859fc8a9e52eb835ba2942ad89959f80f5ce6348c95c75d2e93e2b3a00539eedc0c706d301b7bca5b72f493185bf64ac15b4314198e37a8639e5cc6e0"
```
This transaction can be send through `eth.sendRawTransaction` on a live node or some public service like `https://etherscan.io/pushTx`.