package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/mogza/abstract-go/clients"
)

func main() {
	ctx := context.Background()
	wsClient, err := clients.DialWS("wss://api.testnet.abs.xyz/ws")
	if err != nil {
		panic(err)
	}
	defer wsClient.Close()

	// Example contract addresses
	erc20Addr := common.HexToAddress("0xe4C7fBB0a626ed208021ccabA6Be1566905E2dFc")
	erc721Addr := common.HexToAddress("0x30072084ff8724098cbb65e07f7639ed31af5f66")

	// Watch ERC20 Transfers
	err = wsClient.WatchContractEvent(ctx, erc20Addr, clients.MinimalERC20ABI, "Transfer", nil, func(vLog types.Log) error {
		from := common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		value := new(big.Int).SetBytes(vLog.Data)
		fmt.Printf("[ERC20 Transfer] from %s → %s | value: %s\n", from.Hex(), to.Hex(), value.String())
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("ERC20 watch failed: %w", err))
	}

	// Watch ERC721 Transfers
	err = wsClient.WatchContractEvent(ctx, erc721Addr, clients.MinimalERC721ABI, "Transfer", nil, func(vLog types.Log) error {
		from := common.HexToAddress(vLog.Topics[1].Hex())
		to := common.HexToAddress(vLog.Topics[2].Hex())
		tokenID := new(big.Int).SetBytes(vLog.Topics[3].Bytes())
		fmt.Printf("[ERC721 Transfer] from %s → %s | tokenID: %s\n", from.Hex(), to.Hex(), tokenID.String())
		return nil
	})
	if err != nil {
		panic(fmt.Errorf("ERC721 watch failed: %w", err))
	}

	fmt.Println("🟢 Watching ERC20 + ERC721 transfer events...")
	<-ctx.Done()
	fmt.Println("🛑 Context canceled, stopping watcher.")
}
