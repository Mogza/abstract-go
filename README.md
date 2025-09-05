# abstract-go

`abstract-go` is a Go SDK for building on [Abstract](https://abs.xyz), an EVM-compatible L2.
It wraps the standard [go-ethereum](https://github.com/ethereum/go-ethereum) `ethclient`
and provides Abstract-first naming, helpers, and examples.

## ✨ Features (v0.1 MVP)

- Connect to an Abstract RPC node
- Query balances and nonces
- Get current gas price
- Call contracts
- Send ETH transactions
- Abstract-flavored helpers and examples

## 🚀 Installation

```bash
go get github.com/mogza/abstract-go
```

🛠 Usage
```go
package main

import (
    "context"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/mogza/abstract-go/clients"
)

func main() {
    ctx := context.Background()

    // Connect to Abstract RPC
    client, err := clients.Dial("https://api.mainnet.abs.xyz")
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Example address
    addr := common.HexToAddress("0x0000000000000000000000000000000000000000")

    // Query balance
    balance, err := client.BalanceAt(ctx, addr, nil)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Balance:", balance)
}
```

📂 Project Structure
```bash
.
├── clients/         # Core SDK
│   └── client.go
├── examples/        # Example usage
│   └── transfer.go
├── go.mod
├── LICENSE
└── README.md
```

🔮 Roadmap

v0.2: ERC20 transfers, contract deployment helpers

v0.3: Event logs, filters, subscriptions

Future: Abstract-specific extensions (if/when RPC grows)

🤝 Contributing

PRs and issues are welcome! This SDK is community-driven to help Abstract adoption.

Note: This library builds on top of go-ethereum. Most functionality is identical,
but the goal of abstract-go is to provide a friendly, Abstract-native developer experience
and a future-proof place for Abstract-specific features.
