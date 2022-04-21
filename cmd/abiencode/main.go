package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/eth"
	"log"
	"math/big"
	"time"
)

func main() {
	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("52123842115ecd8af0668615cea556d4afbcd4184394321831c7f80c896f0226")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NONCE: %+v\n", nonce)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("GasPrice: %+v\n", gasPrice)

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress("0x302A600C5a1f81F4d9691147DC6B2f8E77e151F0")
	s, err := eth.NewEth(address, client)
	if err != nil {
		log.Fatal(err)
	}

	tString, err := abi.NewType("string", "", nil)
	if err != nil {
		log.Fatal(err)
	}

	tUint64, err := abi.NewType("uint64", "", nil)
	if err != nil {
		log.Fatal(err)
	}

	args := abi.Arguments{
		{Type: tString},
		{Type: tUint64},
		{Type: tString},
		{Type: tUint64},
	}
	b, err := args.Pack("prices", uint64(time.Now().Unix()), "ETH", uint64(260780000))
	if err != nil {
		log.Fatal(err)
	}

	//hashBytes := crypto.Keccak256Hash(b)
	sign := signatureHash(b)

	tx, err := s.Put(auth, b, sign)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TX")
	log.Println(tx.Value())

}

func signatureHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
