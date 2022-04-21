package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

const rawUrl = "http://localhost:7545"
const hexPrivateKey = "a67a38a5e0f83b2ecda557179845334313921eee4e6c64743dde36a1ab243e60"
const contractAddressHex = "0x7f5199F1D1845Bf7Fd6E2e67e603e3f72a53f531"

func sig32(data []byte) []byte {
	return crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", crypto.Keccak256(data))))
	//[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", crypto.Keccak256Hash([]byte(data)).Bytes()))
}

func signature(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256(crypto.Keccak256Hash([]byte(msg)).Bytes())
}
