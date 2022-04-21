package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/contracts"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
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

	contractAddress := common.HexToAddress(contractAddressHex)

	c, err := contracts.NewAddress(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	msg := []byte("message")

	log.Println(crypto.Keccak256Hash(msg).Bytes())

	sig, err := crypto.Sign(crypto.Keccak256Hash(msg).Bytes(), privateKey)
	if err != nil {
		log.Println("here")
		log.Fatal(err)
	}
	//log.Println(hex.EncodeToString(sig))
	//log.Println(hex.EncodeToString(crypto.Keccak256Hash(msg).Bytes()))
	//
	//log.Println(len(sig))
	//log.Println(hex.EncodeToString(sig[:32]), hex.EncodeToString(sig[32:64]), hex.EncodeToString(sig[64:]))

	log.Println(hex.EncodeToString(msg))
	log.Println(hex.EncodeToString(sig))

	tx, err := c.Source(&bind.CallOpts{
		Pending:     false,
		From:        fromAddress,
		BlockNumber: nil,
		Context:     context.Background(),
	}, msg, sig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tx.Value())
	log.Println(tx.String())
	//s, err := eth.NewEth(address, client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//tString, err := abi.NewType("string", "", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//tUint64, err := abi.NewType("uint64", "", nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//args := abi.Arguments{
	//	{Type: tString},
	//	{Type: tUint64},
	//	{Type: tString},
	//	{Type: tUint64},
	//}
	//b, err := args.Pack("prices", uint64(time.Now().Unix()), "ETH", uint64(260780000))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	////hashBytes := crypto.Keccak256Hash(b)
	//sign := signatureHash(b)
	//
	//tx, err := s.Put(auth, b, sign)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println("TX")
	//log.Println(tx.Value())

}

func signatureHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
