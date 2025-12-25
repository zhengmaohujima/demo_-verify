package homework05

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockInfo 封装区块关键信息（便于测试 & 复用）
type BlockInfo struct {
	Number     uint64
	Hash       common.Hash
	ParentHash common.Hash
	Timestamp  time.Time
	TxCount    int
	GasUsed    uint64
	GasLimit   uint64
	Miner      common.Address
}

// GetBlockInfoByNumber 查询指定区块号的区块信息
func GetBlockInfoByNumber(
	ctx context.Context,
	client *ethclient.Client,
	blockNumber *big.Int,
) (*BlockInfo, error) {

	block, err := client.BlockByNumber(ctx, blockNumber)
	if err != nil {
		return nil, err
	}

	info := &BlockInfo{
		Number:     block.Number().Uint64(),
		Hash:       block.Hash(),
		ParentHash: block.ParentHash(),
		Timestamp:  time.Unix(int64(block.Time()), 0),
		TxCount:    len(block.Transactions()),
		GasUsed:    block.GasUsed(),
		GasLimit:   block.GasLimit(),
		Miner:      block.Coinbase(),
	}

	return info, nil
}
