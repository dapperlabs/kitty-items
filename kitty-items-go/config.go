package main

import (
	"fmt"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
)

type Config struct {
	MinterFlowAddressHex      string          `required:"true"`
	MinterPrivateKeyHex       string          `required:"true"`
	KibblesContractAddressHex string          `required:"true"`
	FlowNode                  string          `default:"localhost:3569"`
	MinterSigAlgoName         string          `default:"ECDSA_P256"`
	MinterHashAlgoName        string          `default:"SHA3_256"`
	MinterAccountKeyIndex     int             `default:"0"`
	FlowEnvironment           FlowEnvironment `default:"flow-testnet"` //testnet, mainnet or emulator

	// These are computed variables based on the env variables above
	MinterPrivateKey             crypto.PrivateKey `ignored:"true"`
	MinterFlowAddress            flow.Address      `ignored:"true"`
	FungibleTokenContractAddress flow.Address      `ignored:"true"`
	KibblesContractAddress       flow.Address      `ignored:"true"`
}

type FlowEnvironment string

const (
	Emulator FlowEnvironment = "Emulator"
	Testnet  FlowEnvironment = "Testnet"
	Mainnet  FlowEnvironment = "Mainnet"
)

func (e *FlowEnvironment) Decode(value string) error {
	switch value {
	case "flow-testnet":
		*e = Testnet
	case "flow-emulator":
		*e = Emulator
	case "flow-mainnet":
		*e = Mainnet
	default:
		return fmt.Errorf("invalid flow environment = %s", value)
	}
	return nil
}

// fungibleContractFor returns the contract addresses for the core interfaces for Fungible Token, as in:
// https://docs.onflow.org/protocol/core-contracts/fungible-token/
func fungibleContractFor(environment FlowEnvironment) string {
	switch environment {
	case Emulator:
		return "ee82856bf20e2aa6"
	case Testnet:
		return "9a0766d93b6608b7"
	case Mainnet:
		return "f233dcee88fe0abe"
	}
	return ""
}

// Compute sanitizes and converts configurations to their proper types for flow
func (c *Config) Compute() (err error) {
	c.MinterFlowAddress = flow.HexToAddress(c.MinterFlowAddressHex)
	c.MinterPrivateKey, err = crypto.DecodePrivateKeyHex(crypto.StringToSignatureAlgorithm(c.MinterSigAlgoName), c.MinterPrivateKeyHex)
	if err != nil {
		return fmt.Errorf("error decrypting private key: %w", err)
	}
	c.FungibleTokenContractAddress = flow.HexToAddress(fungibleContractFor(c.FlowEnvironment))
	c.KibblesContractAddress = flow.HexToAddress(c.KibblesContractAddressHex)
	return nil
}
