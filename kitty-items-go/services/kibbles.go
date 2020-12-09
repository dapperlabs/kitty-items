package services

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"text/template"

	"github.com/dapperlabs/kitty-items-go/templates"
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type KibblesService struct {
	flowService         *FlowService
	mintKibblesTemplate string
}

type KibblesAddresses struct {
	FungibleTokenAddress string
	KibbleTokenAddress   string
}

func NewKibbles(service *FlowService) (*KibblesService, error) {
	mintKibblesTemplate, err := template.New("mintKibbles").Parse(templates.MintKibblesTemplate)
	if err != nil {
		return nil, err
	}

	var minterCode *bytes.Buffer
	if err := mintKibblesTemplate.Execute(minterCode, KibblesAddresses{
		FungibleTokenAddress: service.FungibleContractAddress.String(),
		KibbleTokenAddress:   service.KibblesContractAddress.String(),
	}); err != nil {
		return nil, fmt.Errorf("error mint template code = %w", err)
	}

	return &KibblesService{flowService: service, mintKibblesTemplate: minterCode.String()}, nil
}

// Mint sends a transaction to the Flow blockchain and returns the generated transactionID as a string.
func (k *KibblesService) Mint(ctx context.Context, destinationAddress flow.Address, amountStr string) (string, error) {
	log.Printf("minting kibbles to address=%s", destinationAddress.String())
	sequenceNumber, err := k.flowService.GetMinterAddressSequenceNumber(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting sequence number = %w", err)
	}

	referenceBlock, err := k.flowService.client.GetLatestBlock(ctx, true)
	if err != nil {
		return "", fmt.Errorf("error getting reference block = %w", err)
	}

	tx := flow.NewTransaction().
		SetScript([]byte(k.mintKibblesTemplate)).
		SetProposalKey(k.flowService.minterAddress, k.flowService.minterAccountKey.Index, sequenceNumber).
		SetPayer(k.flowService.minterAddress).
		SetReferenceBlockID(referenceBlock.ID).
		SetGasLimit(100)

	if err := tx.AddArgument(cadence.NewAddress(destinationAddress)); err != nil {
		return "", fmt.Errorf("invalid flow destination address = %s", err)
	}

	amount, err := cadence.NewUFix64(amountStr)
	if err != nil {
		return "", fmt.Errorf("invalid amount = %s", err)
	}

	if err := tx.AddArgument(amount); err != nil {
		return "", err
	}

	txID, err := k.flowService.Send(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("error sending flow transaction = %w", err)
	}

	return txID, nil
}
