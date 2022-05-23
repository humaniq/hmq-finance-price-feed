package storage

import (
	"context"
	"crypto/ecdsa"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/contracts"
)

type PricesContractSetter struct {
	client     *ethclient.Client
	transactor *contracts.PriceDataTransactor
	opts       *bind.TransactOpts
	address    common.Address
}

func NewPricesContractSetter(rawUrl string, chainId int64, contractAddressHex string, privateKeyString string) (*PricesContractSetter, error) {
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
	return &PricesContractSetter{
		client:     client,
		transactor: p,
		opts:       auth,
		address:    clientAddress,
	}, nil
}

func (pc *PricesContractSetter) SavePrices(ctx context.Context, key string, value *price.Asset) error {
	for _, price := range value.Prices {
		if err := pc.SetSymbolPrice(ctx, price.Symbol, price.Currency, price.Price, price.TimeStamp); err != nil {
			return err
		}
	}
	return nil
}
func (pc *PricesContractSetter) SetSymbolPrices(ctx context.Context, symbol string, timeStamp time.Time, prices map[string]float64) error {
	for currency, price := range prices {
		if err := pc.SetSymbolPrice(ctx, symbol, currency, price, timeStamp); err != nil {
			return err
		}
	}
	return nil
}
func (pc *PricesContractSetter) SetSymbolPrice(ctx context.Context, symbol string, currency string, price float64, timeStamp time.Time) error {
	nonce, err := pc.client.PendingNonceAt(ctx, pc.address)
	if err != nil {
		return err
	}
	auth := pc.opts
	auth.Nonce = big.NewInt(int64(nonce))
	if _, err := pc.transactor.PutPrice(auth, symbol, currency, uint64(math.Round(price*1000000)), uint64(timeStamp.Unix())); err != nil {
		return err
	}
	return nil
}

type PricesContractGetter struct {
	client        *ethclient.Client
	caller        *contracts.PriceDataCaller
	clientAddress common.Address
	sourceAddress common.Address
}

func NewPricesContractGetter(rawUrl string, chainId int64, contractAddressHex string, clientAddressHex string, sourceAddressHex string) (*PricesContractGetter, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}
	p, err := contracts.NewPriceDataCaller(common.HexToAddress(contractAddressHex), client)
	if err != nil {
		return nil, err
	}
	return &PricesContractGetter{
		client:        client,
		caller:        p,
		clientAddress: common.HexToAddress(clientAddressHex),
		sourceAddress: common.HexToAddress(sourceAddressHex),
	}, nil
}
func (pc *PricesContractGetter) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*price.Value, error) {
	value, ts, err := pc.caller.GetPrice(
		&bind.CallOpts{
			From:    pc.clientAddress,
			Context: ctx,
		},
		pc.sourceAddress, symbol, currency,
	)
	if err != nil {
		return nil, err
	}
	record := price.Value{
		TimeStamp: time.Unix(int64(ts), 0),
		Source:    "contract",
		Symbol:    symbol,
		Currency:  currency,
		Price:     float64(value / 1000000),
	}
	return &record, nil
}
