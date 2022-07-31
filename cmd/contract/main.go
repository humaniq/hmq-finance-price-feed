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
	"math"
	"math/big"
	"os"
)

const coinAddressHex = "0x2bA64EFB7A4Ec8983E22A49c81fa216AC33f383A"

const contractAddressHex = "0x5641d752955C089Df17D04A3ac1EFD5Ba590fB7b"

const price = 0.07057921

//const chainId = 1337
const chainId = 56

func main() {

	ctx := context.Background()

	//ethClient, err := ethclient.Dial("http://ssh.cubic.local:30101")
	ethClient, err := ethclient.Dial("https://bsc-dataseed1.binance.org/")
	if err != nil {
		log.Fatal(err)
	}

	privateKeyValue := os.Getenv("PRIVATE_KEY")
	log.Println(privateKeyValue)

	privateKey, err := crypto.HexToECDSA(privateKeyValue)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	clientAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	log.Println(clientAddress)

	//clientAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainId))
	if err != nil {
		log.Fatal(err)
	}
	auth.From = clientAddress
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)

	gasPrice, err := ethClient.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	auth.GasPrice = gasPrice
	gp, err := contracts.NewContractsTransactor(common.HexToAddress(contractAddressHex), ethClient)
	if err != nil {
		log.Fatal(err)
	}

	nonce, err := ethClient.PendingNonceAt(ctx, clientAddress)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	if tx, err := gp.SetDirectPrice(auth, common.HexToAddress(coinAddressHex), toBigIntPrice(price, 18)); err != nil {
		log.Fatal(err)
	} else {
		log.Println(tx.Cost())
		log.Println(tx.GasPrice())
		log.Println(tx.Hash())
		log.Println(tx.Value())
	}

	gpo, err := contracts.NewContractsCaller(common.HexToAddress(contractAddressHex), ethClient)
	//gpo, err := contracts.NewContractsCaller(common.HexToAddress("0x5641d752955C089Df17D04A3ac1EFD5Ba590fB7b"), ethClient)
	if err != nil {
		log.Fatal(err)
	}
	prices, err := gpo.AssetPrices(
		&bind.CallOpts{
			Context: ctx,
		},
		common.HexToAddress(coinAddressHex),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("PRICES: %+v", toFloatPrice(prices, 18))
}

func toBigIntPrice(value float64, decimals float64) *big.Int {
	log.Println(math.Pow(10, 18))
	return big.NewInt(
		int64(value * math.Pow(10, decimals)),
	)
}
func toFloatPrice(value *big.Int, decimals float64) float64 {
	return float64(value.Int64()) / math.Pow(10, decimals)
}
