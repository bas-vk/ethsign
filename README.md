## Ethereum signer
Ethsigner is a utility that uses local keystore and can sign transactions using an interactive console.
It is fully compatible with geth keystore.

### Usage
Ethsigner provides the following features:

`eth.accounts` lists accounts in the keystore

`personal.newAccount([password])` create new account, if password is not given you will be prompted

`personal.unlockAccount(<address>, [password], [duration])` unlock account, if password is not given you will be prompted.

`eth.signTransaction({transaction object})` sign transaction

### Examples
#### Start ethsign
```
ethsign --keystore </path/to/keystore>
```

ethsign optionally accepts the `--lightkdf` flag.
With this flag set it costs less resources to protect your keys at the expense of security.
This is useful when ethsign is run on an environment with low resources.

#### Create account
```
> personal.newAccount()
Passphrase: 
Repeat passphrase: 
"0x33178afe528a64aca3a94b80be341216d61112a1"
```

#### List accounts
```
> eth.accounts
["0x33178afe528a64aca3a94b80be341216d61112a1"]
```

#### Unlock account for 1 minute
```
> personal.unlockAccount("0x33178afe528a64aca3a94b80be341216d61112a1", null, 60)
Unlock account 0x33178afe528a64aca3a94b80be341216d61112a1
Passphrase: 
true
```

The second (password) and third (duration is seconds) are optional. If the password isn't supplied ethsign will ask for it. If the duration isn't given the account is unlocked until ethsign is closed.

#### Sign (value) transaction
```
> eth.signTransaction({from: "0x33178afe528a64aca3a94b80be341216d61112a1", to: "0x11fe4f04a5bfda50e155b4289ed69d6fc348333b", value: 1, nonce: 0, gas: 21000, gasPrice: web3.toWei(20, "shannon"), chainId: 3})
  "0xf864808504a817c8008252089411fe4f04a5bfda50e155b4289ed69d6fc348333b018029a0185316f17dd159019f2b54cb19b2f61b364fe0a312de1b6d3f609c3c536ee41da00efb637f55d0c86be966c66c4691f9c2f012feec2b5a8ca9fe5548a52d4037e0"
```

This transaction can be send through `eth.sendRawTransaction` on a live node or some public service like `https://etherscan.io/pushTx`.

Note: the chainId field is optional and default by default to the main net (see for more information https://blog.ethereum.org/2016/11/20/from-morden-to-ropsten/).

