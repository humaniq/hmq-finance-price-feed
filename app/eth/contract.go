package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Contract struct {
	client          *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	publicKey       *ecdsa.PublicKey
	contractAddress common.Address
	gasLimit        uint64
	from            common.Address
}

func NewContract(rawUrl string, contractAddress string, privateKey string, gasLimit uint64) (*Contract, error) {
	client, err := ethclient.Dial(rawUrl)
	if err != nil {
		return nil, err
	}
	privKey, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}
	pubKey, ok := privKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*pubKey)
	return &Contract{
		client:          client,
		privateKey:      privKey,
		publicKey:       pubKey,
		contractAddress: common.HexToAddress(contractAddress),
		gasLimit:        gasLimit,
		from:            fromAddress,
	}, nil
}

func (c *Contract) Put(ctx context.Context, source string, timeStamp time.Time, key string, value float64) error {
	nonce, err := c.client.PendingNonceAt(context.Background(), c.from)
	if err != nil {
		return err
	}

	gasPrice, err := c.client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	auth := bind.NewKeyedTransactor(c.privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0) // in wei
	auth.GasLimit = c.gasLimit // in units
	auth.GasPrice = gasPrice

	instance, err := NewEthTransactor(c.contractAddress, c.client)
	if err != nil {
		return err
	}

	tx, err := instance.Put(auth, []byte("test"), signatureHash([]byte("test")))
	if err != nil {
		return err
	}
	tx.Type()

	return nil
}
func (c *Contract) Get(ctx context.Context) error {
	return nil
}

func signatureHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}
