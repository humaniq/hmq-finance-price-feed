package svc

import (
	"context"
	"crypto/ecdsa"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/humaniq/hmq-finance-price-feed/app/contracts"
)

type ContractPriceSetter struct {
	client     *ethclient.Client
	transactor *contracts.PriceDataTransactor
	opts       *bind.TransactOpts
	address    common.Address
}

func NewContractPriceSetter(rawUrl string, chainId int64, contractAddressHex string, privateKeyString string) (*ContractPriceSetter, error) {
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
	return &ContractPriceSetter{
		client:     client,
		transactor: p,
		opts:       auth,
		address:    clientAddress,
	}, nil
}
func (cps *ContractPriceSetter) SetSymbolPrice(ctx context.Context, price *PriceRecord) error {
	nonce, err := cps.client.PendingNonceAt(ctx, cps.address)
	if err != nil {
		return err
	}
	auth := cps.opts
	auth.Nonce = big.NewInt(int64(nonce))
	if _, err := cps.transactor.PutPrice(auth, price.Symbol, price.Currency, uint64(math.Round(price.Price*1000000)), uint64(price.TimeStamp.Unix())); err != nil {
		return err
	}
	return nil
}

type ContractPriceGetter struct {
	client        *ethclient.Client
	caller        *contracts.PriceDataCaller
	clientAddress common.Address
	sourceAddress common.Address
}

func NewContractPriceGetter(rawUrl string, chainId int64, contractAddressHex string, clientAddressHex string, sourceAddressHex string) (*ContractPriceGetter, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}
	p, err := contracts.NewPriceDataCaller(common.HexToAddress(contractAddressHex), client)
	if err != nil {
		return nil, err
	}
	return &ContractPriceGetter{
		client:        client,
		caller:        p,
		clientAddress: common.HexToAddress(clientAddressHex),
		sourceAddress: common.HexToAddress(sourceAddressHex),
	}, nil
}
func (cpg *ContractPriceGetter) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
	value, ts, err := cpg.caller.GetPrice(
		&bind.CallOpts{
			From:    cpg.clientAddress,
			Context: ctx,
		},
		cpg.sourceAddress, symbol, currency,
	)
	if err != nil {
		return nil, err
	}
	return &PriceRecord{
		Source:        "contract",
		Symbol:        symbol,
		Currency:      currency,
		Price:         float64(value / 1000000),
		PreviousPrice: 0,
		TimeStamp:     time.Unix(int64(ts), 0),
	}, nil
}
