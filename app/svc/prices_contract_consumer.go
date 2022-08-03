package svc

import (
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/contracts"
	"math"
	"math/big"
)

type ContractPricesConsumer struct {
	client     *ethclient.Client
	transactor *contracts.PriceDataTransactor
	opts       *bind.TransactOpts
	address    common.Address
}

func NewContractPricesConsumer(ctx context.Context, rawUrl string, contractAddressHex string, privateKeyString string, chainId int64) (*ContractPricesConsumer, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyString)
	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, err
	}

	clientAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainId))
	if err != nil {
		return nil, err
	}
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000)

	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	auth.GasPrice = gasPrice
	p, err := contracts.NewPriceDataTransactor(common.HexToAddress(contractAddressHex), client)
	if err != nil {
		return nil, err
	}

	return &ContractPricesConsumer{
		client:     client,
		transactor: p,
		opts:       auth,
		address:    clientAddress,
	}, nil
}
func (cpc *ContractPricesConsumer) Consume(ctx context.Context, feed <-chan *FeedItem) error {
	for item := range feed {
		app.Logger().Info(ctx, "PriceContract: %+v", item)
		for _, price := range item.records {
			nonce, err := cpc.client.PendingNonceAt(ctx, cpc.address)
			if err != nil {
				return err
			}
			auth := cpc.opts
			auth.Nonce = big.NewInt(int64(nonce))
			if _, err := cpc.transactor.PutPrice(auth, price.Symbol, price.Currency, uint64(math.Round(price.Price*1000000)), uint64(price.TimeStamp.Unix())); err != nil {
				app.Logger().Error(ctx, "error contract tx: %s", err.Error())
				continue
			}
		}
	}
	return nil
}
