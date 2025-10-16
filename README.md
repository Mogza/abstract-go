# abstract-go
![GoAbstract-removebg-preview](https://github.com/user-attachments/assets/ea6dc78f-9d59-40a4-9139-c9575bc4c451)

`abstract-go` is a Go SDK for building on [Abstract](https://abs.xyz), an EVM-compatible L2.
It wraps the standard [go-ethereum](https://github.com/ethereum/go-ethereum) `ethclient`
and provides Abstract-first naming, helpers, and examples.

## ✨ Features (v1)   
### Wallet & Keys
- Import/export wallets (private key, mnemonic, keystore JSON)
- Message signing: EIP-191, EIP-712 typed data
- Signature recovery & verification
- Deterministic HD wallets for testing/dev

### Transactions Utilities
- ApproveAndTransferERC20
- BatchSendETH
- SafeContractCall
- Robust, thread-safe nonce manager
- Gas estimation helpers (+ buffers, ERC20 & contract calls)
- Auto-fill transaction builder (BuildAndSendTx) with sane defaults

### ERC20 Support
- balanceOf, transfer, approve, allowance, decimals, symbol, name
- Watchers: Transfer & Approval events (real-time)

### ERC721 (NFT) Support
- balanceOf, ownerOf, tokenURI, transferFrom
- WatchERC721Transfers (real-time)
- Simple ERC721 client struct mirroring the ERC20 client

### Unified Event Watching
- WatchContractEvent(contractAddr, abi, eventName, handlerFn)
- Works with JSON ABI and filters

> Note: This library builds on top of go-ethereum. Most functionality is identical, but the goal of abstract-go is to provide a friendly, Abstract-native developer experience and a future-proof place for Abstract-specific features.

## 🚀 Installation

```bash
go get github.com/mogza/abstract-go
```

## 🛠 Usage   

1️⃣ Create/Import a Wallet
```go
wallet, _ := clients.NewWallet()
wallet2, _ := clients.FromPrivateKey("HEX_KEY")
wallet3, _ := clients.NewWalletFromMnemonic("mnemonic words ...", "")
keyJSON, _ := wallet.ExportKeystoreJSON("password")
wallet4, _ := clients.ImportKeystoreJSON(keyJSON, "password")

```

2️⃣ Sign & Verify Messages
```go
message := []byte("Hello Abstract!")
sig, _ := wallet.SignMessageEIP191(message)
valid, _ := clients.VerifySignature(clients.PrefixedHash(message), sig, wallet.Address)
fmt.Println("Valid signature?", valid)

// EIP-712 typed data
typedHash := clients.HashTypedData(myTypedData)
sig2, _ := wallet.SignHash(typedHash)
addr2, _ := clients.RecoverAddressFromSignature(typedHash, sig2)
fmt.Println("Typed data signer:", addr2.Hex())
```

3️⃣ Send ETH
```go
recipient := common.HexToAddress("0xRecipient")
amount := big.NewInt(10000000000000000) // 0.01 ETH
nm := clients.NewNonceManager(client, wallet.Address)

tx, err := wallet.BuildAndSendTx(ctx, client, &recipient, amount, nil, nm)
fmt.Println("Transaction Hash:", tx.Hash().Hex())
```

4️⃣ ERC20 Helpers
```go
token := common.HexToAddress("0xYourERC20TokenAddress")
recipient := common.HexToAddress("0xRecipientAddress")

// Balance
balance, _ := clients.ERC20BalanceOf(ctx, client, token, wallet.Address)
fmt.Println("ERC20 Balance:", balance.String())

// Transfer
amount := big.NewInt(1000)
tx, _ := clients.ERC20Transfer(ctx, wallet, client, token, recipient, amount)
fmt.Println("ERC20 Transfer Tx Hash:", tx.Hash().Hex())

// Approve
tx, _ = clients.ERC20Approve(ctx, wallet, client, token, recipient, amount)
fmt.Println("ERC20 Approve Tx Hash:", tx.Hash().Hex())
```

5️⃣ BatchSendETH
```go
recipients := []common.Address{common.HexToAddress("0xAddr1"), common.HexToAddress("0xAddr2")}
amounts := []*big.Int{big.NewInt(10000000000000000), big.NewInt(20000000000000000)}
txs, _ := wallet.BatchSendETH(ctx, client, recipients, amounts, nm)
for i, tx := range txs {
    fmt.Printf("ETH tx %d: %s\n", i, tx.Hash().Hex())
}
```

6️⃣ Subscribe to New Blocks
```go
headers := make(chan *types.Header)
sub, err := client.SubscribeNewHeads(context.Background(), headers)
if err != nil {
	log.Fatal(err)
}

fmt.Println("📡 Subscribed to new heads")

for {
	select {
	case err := <-sub.Err():
		log.Println("Subscription error:", err)
		return
	case header := <-headers:
		fmt.Println("⛓ New block:", header.Number, "at", time.Unix(int64(header.Time), 0))
	}
}
```


See other examples in `examples/`.


## 📂 Project Structure
```bash
.
├── clients
│   ├── client.go
│   ├── erc20.go
│   ├── erc721.go
│   ├── nonce.go
│   ├── subscription.go
│   ├── wallet.go
│   └── wallet_utils.go
├── examples
│   ├── client.go
│   ├── erc20.go
│   ├── erc20Watchers.go
│   ├── erc721.go
│   ├── global.go
│   ├── subLogs.go
│   ├── subManager.go
│   ├── subNewHeads.go
│   ├── subPendingTxs.go
│   ├── transfer.go
│   ├── wallet_deterministic.go
│   ├── wallet.go
│   ├── wallet_keystore.go
│   ├── wallet_mnemonic.go
│   ├── wallet_sign.go
│   └── wallet_verify.go
├── go.mod
├── go.sum
├── LICENSE
└── README.md


```

## 🔮 Roadmap   
- Retry & Resilience: polling with retry/backoff, detect reorgs
- More ERC721 helpers & events
- Additional Abstract-specific utilities

  
## 🤝 Contributing    
PRs and issues are welcome! This SDK is community-driven to help Abstract adoption and provide a Go-first developer experience.    
