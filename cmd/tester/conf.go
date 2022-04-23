package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
)

const rawUrl = "http://ganache.cubic.local:30101"
const hexPrivateKey = "0807fc95777ab8046a708d3947aa99f7c6ebf86fe555bbf0f6d35d37885bf7fc"
const contractAddressHex = "0xdde5E7c605473C4F8E46DC5c294C26B0428b184a"

func sig32(data []byte) []byte {
	return crypto.Keccak256([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", crypto.Keccak256(data))))
	//[]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n32%s", crypto.Keccak256Hash([]byte(data)).Bytes()))
}

func signature(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256(crypto.Keccak256Hash([]byte(msg)).Bytes())
}
