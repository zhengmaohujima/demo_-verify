package homework05

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"homework05/contract/counter"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

// SendETHTransaction 发送一笔 ETH 转账交易（Sepolia）
func SendETHTransaction(
	rpcURL string,
	privateKeyHex string,
	toAddress string,
	amountWei *big.Int,
) (common.Hash, error) {

	// 1. 连接 Sepolia
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return common.Hash{}, err
	}
	defer client.Close()

	// 2. 解析私钥
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return common.Hash{}, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return common.Hash{}, errors.New("invalid public key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// 3. 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return common.Hash{}, err
	}

	// 4. 获取 gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return common.Hash{}, err
	}

	// 5. 构造交易
	to := common.HexToAddress(toAddress)
	var data []byte
	gasLimit := uint64(21000)

	tx := types.NewTransaction(
		nonce,
		to,
		amountWei,
		gasLimit,
		gasPrice,
		data,
	)

	// 6. 签名交易（Sepolia chainId = 11155111）
	chainID := big.NewInt(11155111)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return common.Hash{}, err
	}

	// 7. 发送交易
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, err
	}

	return signedTx.Hash(), nil
}

// CallContract 调用计数器合约的 Increment 方法，并返回最新值
func CallContract(rpcURL, contractAddr, privateKeyHex string, chainID int64) (*big.Int, string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect RPC: %w", err)
	}

	addr := common.HexToAddress(contractAddr)
	ctr, err := counter.NewCounter(addr, client)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load contract: %w", err)
	}

	// 私钥处理
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, "", fmt.Errorf("invalid private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		return nil, "", fmt.Errorf("failed to create transactor: %w", err)
	}

	// 调用 Increment
	tx, err := ctr.Increment(auth)
	if err != nil {
		return nil, "", fmt.Errorf("failed to call Increment: %w", err)
	}

	// 读取最新 count
	value, err := ctr.Get(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get count: %w", err)
	}

	return value, tx.Hash().Hex(), nil
}
