package prices

//
//import (
//	"context"
//	"crypto/ecdsa"
//	"fmt"
//	"github.com/ethereum/go-ethereum/accounts/abi/bind"
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/crypto"
//	"github.com/ethereum/go-ethereum/ethclient"
//	"github.com/humaniq/hmq-finance-price-feed/pkg/ethereum/contracts"
//	"log"
//	"math/big"
//	"time"
//
//	"github.com/humaniq/hmq-finance-price-feed/app"
//	"github.com/humaniq/hmq-finance-price-feed/app/config"
//	"github.com/humaniq/hmq-finance-price-feed/app/price"
//	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
//)
//
//type PriceOracleState struct {
//	Value            float64
//	TimeStamp        time.Time
//	SymbolAddressHex string
//	Decimals         float64
//	Threshold        float64
//	ThresholdValue   float64
//}
//
//type PriceOracleFilter struct {
//	state map[string]PriceOracleState
//	write func(ctx context.Context, value PriceOracleState) (string, error)
//}
//
//func NewPriceOracle(cfg config.PriceOracleContract, networks map[string]config.EthNetwork) (*PriceOracleFilter, error) {
//	network, found := networks[cfg.NetworkUid]
//	if !found {
//		return nil, fmt.Errorf("%w: network %s not found", app.ErrValueInvalid, cfg.NetworkUid)
//	}
//	state := make(map[string]PriceOracleState)
//	for _, value := range cfg.Tokens {
//		key := fmt.Sprintf("%s-%s", value.Symbol, value.Currency)
//		networkSymbol, found := network.Symbols[value.Symbol]
//		if !found {
//			return nil, fmt.Errorf("%w: symbol %s not found", app.ErrValueInvalid, value.Symbol)
//		}
//		state[key] = PriceOracleState{
//			Value:            0,
//			TimeStamp:        time.Now(),
//			SymbolAddressHex: networkSymbol.AddressHex,
//			Decimals:         networkSymbol.Decimals,
//			Threshold:        value.PercentThreshold,
//		}
//	}
//	return &PriceOracleFilter{
//		state: state,
//		write: contractWriterFunc(
//			network,
//			cfg,
//		),
//	}, nil
//}
//
//func (pof *PriceOracleFilter) Work(ctx context.Context, values []price.Value) {
//	for _, value := range values {
//		key := fmt.Sprintf("%s-%s", value.Symbol, value.Currency)
//		current, found := pof.state[key]
//		if !found {
//			continue
//		}
//		diff := current.Value - value.Price
//		if diff < 0 {
//			diff = diff * (-1)
//		}
//		if diff < current.ThresholdValue {
//			continue
//		}
//		current.TimeStamp = time.Now()
//		current.Value = value.Price
//		current.ThresholdValue = current.Value * current.Threshold / 100
//		tx, err := pof.write(ctx, current)
//		if err != nil {
//			logger.Error(ctx, "TX ERROR %s", err)
//			continue
//		}
//		pof.state[key] = current
//		logger.Info(ctx, "Written Value: %+v => %s", current, tx)
//	}
//}
//
//func contractWriterFunc(
//	network config.EthNetwork,
//	contract config.PriceOracleContract,
//) func(ctx context.Context, value PriceOracleState) (string, error) {
//	privateKey, err := crypto.HexToECDSA(contract.ClientPrivateKey)
//	if err != nil {
//		log.Fatal(err)
//	}
//	publicKey := privateKey.Public()
//	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
//	if !ok {
//		log.Fatal(err)
//	}
//
//	clientAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
//	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(network.ChainId))
//	if err != nil {
//		log.Fatal(err)
//	}
//	auth.From = clientAddress
//	auth.Value = big.NewInt(0)
//	auth.GasLimit = uint64(300000)
//
//	ethClient, err := ethclient.Dial("https://bsc-dataseed1.binance.org/")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	gp, err := contracts.NewContractsTransactor(common.HexToAddress(contract.ContractAddressHex), ethClient)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return func(ctx context.Context, value PriceOracleState) (string, error) {
//		log.Printf("CONTRACT WRITE: %+v", value)
//		gasPrice, err := ethClient.SuggestGasPrice(context.Background())
//		if err != nil {
//			return "", err
//		}
//		auth.GasPrice = gasPrice
//		nonce, err := ethClient.PendingNonceAt(ctx, clientAddress)
//		if err != nil {
//			return "", err
//		}
//		auth.Nonce = big.NewInt(int64(nonce))
//
//		tx, err := gp.SetDirectPrice(auth, common.HexToAddress(value.SymbolAddressHex), toBigIntPrice(value.Value, value.Decimals))
//		if err != nil {
//			return "", err
//		}
//
//		return tx.Hash().String(), nil
//	}
//}
