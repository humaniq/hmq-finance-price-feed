package main

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/contracts"
	"log"
	"math/big"
)

func main() {

	ctx := context.Background()

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

	c, err := contracts.NewPriceData(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	//tx, err := c.PutPrice(auth, "BTC", "USD", uint64(3420000000), uint64(time.Now().Unix()))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(tx.Value())
	//
	//auth.Nonce = big.NewInt(int64(nonce + 1))
	//tx, err = c.PutEth(auth, "BTC", uint64(270000000), uint64(time.Now().Unix()))
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(tx.Value())

	value, ts, err := c.GetPrice(&bind.CallOpts{
		Pending: false,
		From:    fromAddress,
		Context: ctx,
	}, "BTC", "ETH")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(value, ts)

	value, ts, err = c.GetEth(&bind.CallOpts{
		Pending: false,
		From:    fromAddress,
		Context: ctx,
	}, "BTC")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(value, ts)

}
