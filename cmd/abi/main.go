package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/eth"
	"log"
	"math/big"
)

func main() {

	tPair, err := abi.NewType("")

	tuple, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "kind", Type: "string"},
		{Name: "timestamp", Type: "uint64"},
		{Name: "key", Type: "string"},
		{Name: "value", Type: "uint64"},
	})

	if err != nil {
		log.Fatal(err)
	}

	tArg := abi.Argument{
		Name:    "message",
		Type:    tuple,
		Indexed: false,
	}

	//bts, err := json.Marshal(&struct {
	//	Kind      string
	//	Timestamp uint64
	//	Key       string
	//	Value     uint64
	//}{
	//	Kind:      "p",
	//	Timestamp: uint64(time.Now().Unix()),
	//	Key:       "ETH",
	//	Value:     42,
	//})
	//bts, err := json.Marshal([]interface{}{
	//	"pices",
	//	time.Now().Unix(),
	//	"ETH",
	//	64,
	//})
	b, _ := abi.NewType("bytes", "", nil)
	log.Println(b)
	args := abi.Arguments{
		tArg,
		//{Name: "signature", Type: b},
	}
	bts, err := args.Pack(struct {
		Kind      string
		Timestamp uint64
		Key       string
		Value     uint64
	}{
		Kind:      "prices",
		Timestamp: 42,
		Key:       "ETH",
		Value:     42,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(bts)
	log.Println(hex.EncodeToString(bts))

	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("9d6dde467955b4295cc9faf10d3453a920d55aac5618474257f2df1b00383afe")
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

	address := common.HexToAddress("0x95120502833079cfd7A5346dF3e67dbFdaDc1975")
	s, err := eth.NewEth(address, client)
	if err != nil {
		log.Fatal(err)
	}

	signature, err := crypto.Sign(crypto.Keccak256Hash([]byte("put(bytes,bytes,string)")).Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := s.Put(auth, bts, signature)
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
