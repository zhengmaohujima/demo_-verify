package homework05

import (
	"context"
	"homework05"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGetBlockInfoByNumber(t *testing.T) {
	// 1. 连接 Sepolia
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/xxxxxxxxxxxx")
	if err != nil {
		t.Fatalf("failed to connect to sepolia: %v", err)
	}
	defer client.Close()

	// 2. 选择一个已存在的区块号（不用太新）
	blockNumber := big.NewInt(6000000)

	// 3. 调用方法
	info, err := homework05.GetBlockInfoByNumber(context.Background(), client, blockNumber)
	if err != nil {
		t.Fatalf("GetBlockInfoByNumber error: %v", err)
	}

	// 4. 断言（Assertions）
	if info.Number != blockNumber.Uint64() {
		t.Errorf("expected block number %d, got %d", blockNumber.Uint64(), info.Number)
	}

	if info.Hash.Hex() == "" {
		t.Error("block hash should not be empty")
	}

	if info.TxCount < 0 {
		t.Error("tx count should be >= 0")
	}

	if info.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}

	t.Logf("Block %d | tx=%d | hash=%s",
		info.Number,
		info.TxCount,
		info.Hash.Hex(),
	)
}
