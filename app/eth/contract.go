package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"

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

func (c *Contract) Set(ctx context.Context) error {

}
func (c *Contract) Get(ctx context.Context) error {

}
