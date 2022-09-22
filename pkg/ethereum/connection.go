package ethereum

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ConnectionConfig struct {
	RawUrl     string
	ChainID    int64
	PrivateKey string
	AddressHex string
}

type TransactConnection struct {
	ethClient *ethclient.Client
	auth      *bind.TransactOpts
	txMx      sync.Mutex
}

func NewTransactConnection(rawUrl string, chainId int64, clientPrivateKey string, gasLimit uint64) (*TransactConnection, error) {
	privateKey, err := crypto.HexToECDSA(clientPrivateKey)
	if err != nil {
		return nil, err
	}
	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid publicKey format")
	}
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainId))
	if err != nil {
		return nil, err
	}
	auth.From = crypto.PubkeyToAddress(*publicKeyECDSA)
	auth.Value = big.NewInt(0)
	auth.GasLimit = gasLimit

	conn, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}

	return &TransactConnection{
		ethClient: conn,
		auth:      auth,
	}, nil
}

func (tc *TransactConnection) Client() *ethclient.Client {
	return tc.ethClient
}
func (tc *TransactConnection) TxBegin(ctx context.Context) (*bind.TransactOpts, error) {
	tc.txMx.Lock()
	gasPrice, err := tc.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	nonce, err := tc.Nonce(ctx)
	if err != nil {
		return nil, err
	}
	auth := tc.auth
	auth.GasPrice = gasPrice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 0
	return auth, nil
}
func (tc *TransactConnection) TxEnd() {
	tc.txMx.Unlock()
}
func (tc *TransactConnection) TxWait(ctx context.Context, tx *types.Transaction) {
	bind.WaitMined(ctx, tc.Client(), tx)
}
func (tc *TransactConnection) GasPrice(ctx context.Context) (*big.Int, error) {
	return tc.Client().SuggestGasPrice(ctx)
}
func (tc *TransactConnection) Nonce(ctx context.Context) (uint64, error) {
	return tc.Client().PendingNonceAt(ctx, tc.auth.From)
}
