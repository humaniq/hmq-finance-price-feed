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
	acon "github.com/humaniq/hmq-finance-price-feed/app/contracts/address"
	"log"
	"math/big"
	"time"
)

func signatureHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func main() {
	client, err := ethclient.Dial("http://localhost:7545")
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("ce51202b43d035651a354a509e9e4c9250bf1027e7452622e1c91f2326b1c1ca")
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

	address := common.HexToAddress("0xf0A52CBA50E69129D3b47EC39489A19695e88e21")
	s, err := acon.NewAddress(address, client)
	if err != nil {
		log.Fatal(err)
	}

	m := []byte(fmt.Sprintf("ADDRESS %s UPDATE ADDRESS TIMESTAMP %d", "0x7Fa68118f03d097bCE3F1544a80ee20841FeB56a", time.Now().Unix()))
	sig := signatureHash(m)
	log.Println(hex.EncodeToString(sig))
	tx, err := s.Set(auth, m, sig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tx.Value())
	log.Println(tx.Data())

	v, err := s.Get(&bind.CallOpts{
		Pending:     false,
		From:        fromAddress,
		BlockNumber: nil,
		Context:     nil,
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(v)

}
