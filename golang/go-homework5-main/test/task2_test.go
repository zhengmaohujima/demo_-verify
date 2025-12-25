package homework05

import (
	"fmt"
	"homework05"
	"math/big"
	"os"
	"strings"
	"testing"
)

func TestSendETHTransaction(t *testing.T) {

	// ⚠️ 使用环境变量保存私钥
	privateKeyHex := os.Getenv("PRIVATE_KEY_test001")
	privateKey := strings.TrimPrefix(privateKeyHex, "0x")

	if privateKey == "" {
		t.Fatal("SEPOLIA_PRIVATE_KEY not set")
	}
	INFURA_KEY_001 := os.Getenv("INFURA_KEY_001")

	rpcURL := "https://sepolia.infura.io/v3/" + INFURA_KEY_001
	toAddress := "0xbfF5A8167C1394B3BE82837716Fcee71E84939C1"

	// 0.001 ETH
	amountWei := new(big.Int)
	amountWei.SetString("1000000000000000", 10)

	txHash, err := homework05.SendETHTransaction(
		rpcURL,
		privateKey,
		toAddress,
		amountWei,
	)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Transaction hash:", txHash.Hex())
}

func TestCall(t *testing.T) {
	rpcURL := "https://sepolia.infura.io/v3/" + os.Getenv("INFURA_KEY_001")
	contractAddr := "0xaC7d68AE475b14d0F0067E64e2862fBBb92D47d8"
	privateKey := os.Getenv("PRIVATE_KEY_test001")
	chainID := int64(11155111) // Sepolia

	value, txHash, err := homework05.CallContract(rpcURL, contractAddr, privateKey, chainID)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("Transaction Hash:", txHash)
	fmt.Println("Current count:", value)
}
