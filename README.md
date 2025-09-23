# abstract-go
![GoAbstract-removebg-preview](https://github.com/user-attachments/assets/ea6dc78f-9d59-40a4-9139-c9575bc4c451)

`abstract-go` is a Go SDK for building on [Abstract](https://abs.xyz), an EVM-compatible L2.
It wraps the standard [go-ethereum](https://github.com/ethereum/go-ethereum) `ethclient`
and provides Abstract-first naming, helpers, and examples.

## ✨ Features (v0.3)

- Connect to an Abstract RPC node
- Query ETH balances, nonces, gas prices
- Call contracts & send ETH transactions
- ERC20 support: balance, transfer, approve, allowance, decimals, symbol, name
- WebSocket handling (DialWS)
- Subscriptions: NewHeads, Logs, PendingTxs
- Multi-Subscription Manager
- ERC20 Watchers: real-time Transfer & Approval events

> Note: This library builds on top of go-ethereum. Most functionality is identical, but the goal of abstract-go is to provide a friendly, Abstract-native developer experience and a future-proof place for Abstract-specific features.

## 🚀 Installation

```bash
go get github.com/mogza/abstract-go
```

## 🛠 Usage   

1️⃣ Create a Wallet
```go
package main

import (
    "fmt"
    "log"

    "github.com/mogza/abstract-go/clients"
)

func main() {
    wallet, err := clients.NewWallet()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Wallet Address:", wallet.Address.Hex())
    fmt.Println("Private Key:", wallet.PrivateKey.D.Text(16)) // Keep secret!
}
```

2️⃣ Connect and Query ETH Balance
```go
ctx := context.Background()
client, _ := clients.DialHTTP("https://api.testnet.abs.xyz")
defer client.Close()

address := common.HexToAddress("0x0000000000000000000000000000000000000000")
balance, _ := client.BalanceAt(ctx, address)
fmt.Println("ETH Balance:", balance.String())
```

3️⃣ Send ETH
```go
recipient := common.HexToAddress("0xRecipientAddress")
amount := big.NewInt(0).SetString("10000000000000000", 10) // 0.01 ETH

err := wallet.SendETH(ctx, client, recipient, amount)
if err != nil {
    log.Fatal(err)
}

fmt.Println("✅ ETH sent!")
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

5️⃣ Subscribe to New Blocks
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
│   ├── client.go
│   ├── erc20.go
│   ├── subscription.go
│   └── wallet.go
├── examples
│   ├── client.go
│   ├── erc20.go
│   ├── erc20Watchers.go
│   ├── global.go
│   ├── subLogs.go
│   ├── subManager.go
│   ├── subNewHeads.go
│   ├── subPendingTxs.go
│   ├── transfer.go
│   └── wallet.go
├── go.mod
├── go.sum
├── LICENSE
└── README.md

```

## 🔮 Roadmap    
v0.4:  
- Wallet management (Export private key, Sign arbitrary messages (EIP-191), Sign typed data (EIP-712))    
- Transaction utils (Gas estimation handler, Nonce manager)
     
Upcoming:  
- ERC721 (NFT) Support  
- Unified Event Watching  
- Transaction Builders  
- Tooling (CLI, ...)  

## 🤝 Contributing    
PRs and issues are welcome! This SDK is community-driven to help Abstract adoption and provide a Go-first developer experience.    
