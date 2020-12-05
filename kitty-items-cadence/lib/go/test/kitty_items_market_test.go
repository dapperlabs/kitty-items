package test

import (
	"strings"
	"testing"

	"github.com/onflow/cadence"
	emulator "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-go-sdk"
	sdk "github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"
	"github.com/stretchr/testify/assert"
)

const (
	kittyItemsMarketRootPath              = "../../../cadence/KittyItemsMarket"
	kittyItemsMarketKittyItemsMarketPath  = kittyItemsMarketRootPath + "/contracts/KittyItemsMarket.cdc"
	kittyItemsMarketSetupAccountPath      = kittyItemsMarketRootPath + "/transactions/setup_account.cdc"
	kittyItemsMarketListItemOfferPath     = kittyItemsMarketRootPath + "/transactions/list_item_offer.cdc"
	kittyItemsMarketPurchaseItemOfferPath = kittyItemsMarketRootPath + "/transactions/purchase_item_offer.cdc"
)

const (
	typeID1337 = 1337
)

type TestContractsInfo struct {
	FTAddr                 flow.Address
	KibbleAddr             flow.Address
	KibbleSigner           crypto.Signer
	NFTAddr                flow.Address
	KittyItemsAddr         flow.Address
	KittyItemsSigner       crypto.Signer
	KittyItemsMarketAddr   flow.Address
	KittyItemsMarketSigner crypto.Signer
}

func KittyItemsMarketDeployContracts(b *emulator.Blockchain, t *testing.T) TestContractsInfo {
	accountKeys := test.AccountKeyGenerator()

	ftAddr, kibbleAddr, kibbleSigner := KibbleDeployContracts(b, t)
	nftAddr, kittyItemsAddr, kittyItemsSigner := KittyItemsDeployContracts(b, t)

	// Should be able to deploy a contract as a new account with one key.
	kittyItemsMarketAccountKey, kittyItemsMarketSigner := accountKeys.NewWithSigner()
	kittyItemsMarketCode := loadKittyItemsMarket(
		ftAddr.String(),
		nftAddr.String(),
		kibbleAddr.String(),
		kittyItemsAddr.String(),
	)
	kittyItemsMarketAddr, err := b.CreateAccount(
		[]*flow.AccountKey{kittyItemsMarketAccountKey},
		[]sdktemplates.Contract{
			{
				Name:   "KittyItemsMarket",
				Source: string(kittyItemsMarketCode),
			},
		})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	return TestContractsInfo{
		ftAddr,
		kibbleAddr,
		kibbleSigner,
		nftAddr,
		kittyItemsAddr,
		kittyItemsSigner,
		kittyItemsMarketAddr,
		kittyItemsMarketSigner,
	}
}

func KittyItemsMarketSetupAccount(b *emulator.Blockchain, t *testing.T, userAddress sdk.Address, userSigner crypto.Signer, contracts TestContractsInfo) {
	tx := flow.NewTransaction().
		SetScript(kittyItemsMarketGenerateSetupAccountScript(contracts.KittyItemsMarketAddr.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		false,
	)
}

// Create a new account with the Kibble and KittyItems resources set up
func KittyItemsMarketCreatePurchaserAccount(b *emulator.Blockchain, t *testing.T, contracts TestContractsInfo) (sdk.Address, crypto.Signer) {
	userAddress, userSigner, _ := createAccount(t, b)
	KibbleSetupAccount(t, b, userAddress, userSigner, contracts.FTAddr, contracts.KibbleAddr)
	KittyItemsSetupAccount(t, b, userAddress, userSigner, contracts.NFTAddr, contracts.KittyItemsAddr)
	return userAddress, userSigner
}

// Create a new account with the Kibble, KittyItems, and KittyItemsMarket resources set up
func KittyItemsMarketCreateAccount(b *emulator.Blockchain, t *testing.T, contracts TestContractsInfo) (sdk.Address, crypto.Signer) {
	userAddress, userSigner := KittyItemsMarketCreatePurchaserAccount(b, t, contracts)
	KittyItemsMarketSetupAccount(b, t, userAddress, userSigner, contracts)
	return userAddress, userSigner
}

func KittyItemsMarketListItem(b *emulator.Blockchain, t *testing.T, contracts TestContractsInfo, userAddress sdk.Address, userSigner crypto.Signer, tokenID uint64, price string, shouldFail bool) {
	tx := flow.NewTransaction().
		SetScript(kittyItemsMarketGenerateListItemOfferScript(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)
	tx.AddArgument(cadence.NewUInt64(tokenID))
	tx.AddArgument(CadenceUFix64(price))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func KittyItemsMarketPurchaseItem(
	b *emulator.Blockchain,
	t *testing.T,
	contracts TestContractsInfo,
	userAddress sdk.Address,
	userSigner crypto.Signer,
	marketCollectionAddress sdk.Address,
	tokenID uint64,
	shouldFail bool,
) {
	tx := flow.NewTransaction().
		SetScript(kittyItemsMarketGeneratePurchaseItemOfferScript(contracts)).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(userAddress)
	tx.AddArgument(cadence.NewUInt64(tokenID))
	tx.AddArgument(cadence.NewAddress(marketCollectionAddress))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, userAddress},
		[]crypto.Signer{b.ServiceKey().Signer(), userSigner},
		shouldFail,
	)
}

func TestKittyItemsMarketDeployContracts(t *testing.T) {
	b := newEmulator()
	KittyItemsMarketDeployContracts(b, t)
}

func TestKittyItemsMarketSetupAccount(t *testing.T) {
	b := newEmulator()

	contracts := KittyItemsMarketDeployContracts(b, t)

	t.Run("Should be able to create an empty Collection", func(t *testing.T) {
		userAddress, userSigner, _ := createAccount(t, b)
		KittyItemsMarketSetupAccount(b, t, userAddress, userSigner, contracts)
	})
}

func TestKittyItemsMarketCreateSaleOffer(t *testing.T) {
	b := newEmulator()

	contracts := KittyItemsMarketDeployContracts(b, t)

	t.Run("Should be able to create a sale offer, list it, and have it accepted", func(t *testing.T) {
		tokenToList := uint64(0)
		tokenPrice := "1.11"
		userAddress, userSigner := KittyItemsMarketCreateAccount(b, t, contracts)
		// Contract mints item
		KittyItemsMintItem(
			b,
			t,
			contracts.NFTAddr,
			contracts.KittyItemsAddr,
			contracts.KittyItemsSigner,
			typeID1337,
		)
		// Contract transfers item to another seller account (we don't need to do this)
		KittyItemsTransferItem(
			b,
			t,
			contracts.NFTAddr,
			contracts.KittyItemsAddr,
			contracts.KittyItemsSigner,
			tokenToList,
			userAddress,
			false,
		)
		// Other seller account lists the item
		KittyItemsMarketListItem(
			b,
			t,
			contracts,
			userAddress,
			userSigner,
			tokenToList,
			tokenPrice,
			false,
		)
		buyerAddress, buyerSigner := KittyItemsMarketCreatePurchaserAccount(b, t, contracts)
		// Fund the purchase
		KibbleMint(
			t,
			b,
			contracts.FTAddr,
			contracts.KibbleAddr,
			contracts.KibbleSigner,
			buyerAddress,
			"100.0",
		)
		KittyItemsMarketPurchaseItem(
			b,
			t,
			contracts,
			buyerAddress,
			buyerSigner,
			userAddress,
			tokenToList,
			false,
		)
	})
}

func replaceKittyItemsMarketAddressPlaceholders(codeBytes []byte, contracts TestContractsInfo) []byte {
	code := string(codeBytes)

	code = strings.ReplaceAll(code, ftAddressPlaceholder, "0x"+contracts.FTAddr.String())
	code = strings.ReplaceAll(code, kibbleAddressPlaceHolder, "0x"+contracts.KibbleAddr.String())
	code = strings.ReplaceAll(code, nftAddressPlaceholder, "0x"+contracts.NFTAddr.String())
	code = strings.ReplaceAll(code, kittyItemsAddressPlaceHolder, "0x"+contracts.KittyItemsAddr.String())
	code = strings.ReplaceAll(code, kittyItemsMarketPlaceholder, "0x"+contracts.KittyItemsMarketAddr.String())

	return []byte(code)
}

func loadKittyItemsMarket(ftAddr, nftAddr, kibbleAddr, kittyItemsAddr string) []byte {
	code := string(readFile(kittyItemsMarketKittyItemsMarketPath))

	code = strings.ReplaceAll(code, ftAddressPlaceholder, "0x"+ftAddr)
	code = strings.ReplaceAll(code, kibbleAddressPlaceHolder, "0x"+kibbleAddr)
	code = strings.ReplaceAll(code, nftAddressPlaceholder, "0x"+nftAddr)
	code = strings.ReplaceAll(code, kittyItemsAddressPlaceHolder, "0x"+kittyItemsAddr)

	return []byte(code)
}

func kittyItemsMarketGenerateSetupAccountScript(kittyItemsMarketAddr string) []byte {
	code := string(readFile(kittyItemsMarketSetupAccountPath))

	code = strings.ReplaceAll(code, kittyItemsMarketPlaceholder, "0x"+kittyItemsMarketAddr)

	return []byte(code)
}

func kittyItemsMarketGenerateListItemOfferScript(contracts TestContractsInfo) []byte {
	return replaceKittyItemsMarketAddressPlaceholders(
		readFile(kittyItemsMarketListItemOfferPath),
		contracts,
	)
}

func kittyItemsMarketGeneratePurchaseItemOfferScript(contracts TestContractsInfo) []byte {
	return replaceKittyItemsMarketAddressPlaceholders(
		readFile(kittyItemsMarketPurchaseItemOfferPath),
		contracts,
	)
}
