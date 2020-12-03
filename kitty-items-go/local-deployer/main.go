package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"io/ioutil"

	"encoding/hex"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/examples"
	"github.com/onflow/flow-go-sdk/templates"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()
	flowClient, err := client.New("127.0.0.1:3569", grpc.WithInsecure())
	handle(err)

	serviceAcctAddr, serviceAcctKey, serviceSigner := examples.ServiceAccount(flowClient)

	pubKeyHex, privKeyHex := generateKeys("ECDSA_P256")
	fmt.Printf("Public Key = %s", pubKeyHex)
	fmt.Printf("Private Key = %s", privKeyHex)

	myPrivateKey, err := crypto.DecodePrivateKeyHex(crypto.StringToSignatureAlgorithm("ECDSA_P256"), privKeyHex)
	handle(err)

	contractCode, err := readFile("./kibbles-nft.cdc")
	handle(err)

	myAcctKey := flow.NewAccountKey().
		FromPrivateKey(myPrivateKey).
		SetHashAlgo(crypto.SHA3_256).
		SetWeight(flow.AccountKeyWeightThreshold)

	referenceBlockID := examples.GetReferenceBlockId(flowClient)
	createAccountTx := templates.CreateAccount([]*flow.AccountKey{myAcctKey}, []templates.Contract{
		{"KibblesLocal", contractCode},
	}, serviceAcctAddr)
	createAccountTx.SetProposalKey(
		serviceAcctAddr,
		serviceAcctKey.Index,
		serviceAcctKey.SequenceNumber,
	)
	createAccountTx.SetReferenceBlockID(referenceBlockID)
	createAccountTx.SetPayer(serviceAcctAddr)

	err = createAccountTx.SignEnvelope(serviceAcctAddr, serviceAcctKey.Index, serviceSigner)
	handle(err)

	err = flowClient.SendTransaction(ctx, *createAccountTx)
	handle(err)

	accountCreationTxRes := examples.WaitForSeal(ctx, flowClient, createAccountTx.ID())
	examples.Handle(accountCreationTxRes.Error)

	// Successful Tx, increment sequence number
	serviceAcctKey.SequenceNumber++

	var myAddress flow.Address

	for _, event := range accountCreationTxRes.Events {
		if event.Type == flow.EventAccountCreated {
			accountCreatedEvent := flow.AccountCreatedEvent(event)
			myAddress = accountCreatedEvent.Address()
		}
	}

	fmt.Println("My Address:", myAddress.Hex())
}

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(contents), nil
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

func generateKeys(sigAlgoName string) (string, string) {
	seed := make([]byte, crypto.MinSeedLength)

	// [2] create random seed
	_, err := rand.Read(seed)

	if err != nil {
		panic(err)
	}

	// [3]
	sigAlgo := crypto.StringToSignatureAlgorithm(sigAlgoName)
	privateKey, err := crypto.GeneratePrivateKey(sigAlgo, seed)
	if err != nil {
		panic(err)
	}

	// [4]
	publicKey := privateKey.PublicKey()

	pubKeyHex := hex.EncodeToString(publicKey.Encode())
	privKeyHex := hex.EncodeToString(privateKey.Encode())

	return pubKeyHex, privKeyHex
}
