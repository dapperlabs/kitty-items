package test

import (
	"strings"
	"testing"

	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/flow-go-sdk/crypto"
	sdktemplates "github.com/onflow/flow-go-sdk/templates"
	"github.com/onflow/flow-go-sdk/test"

	"github.com/onflow/flow-nft/lib/go/contracts"

	"github.com/stretchr/testify/assert"

	"github.com/onflow/flow-go-sdk"

	nft_contracts "github.com/onflow/flow-nft/lib/go/contracts"
)

const (
	kittyItemsRootPath                   = "../../../cadence/kittyItems"
	kittyItemsKittyItemsPath             = kittyItemsRootPath + "/contracts/KittyItems.cdc"
	kittyItemsSetupAccountPath           = kittyItemsRootPath + "/transactions/setup_account.cdc"
	kittyItemsMintKittyItemPath          = kittyItemsRootPath + "/transactions/mint_kitty_item.cdc"
	kittyItemsInspectKittyItemSupplyPath = kittyItemsRootPath + "/scripts/read_kitty_items_supply.cdc"
	kittyItemsInspectCollectionLenPath   = kittyItemsRootPath + "/scripts/read_collection_length.cdc"
	kittyItemsInspectCollectionIdsPath   = kittyItemsRootPath + "/scripts/read_collection_ids.cdc"

	nftAddressPlaceholder        = "0xNONFUNGIBLETOKEN"
	kittyItemsAddressPlaceHolder = "0xKITTYITEMS"

	typeID1 = 1000
	typeID2 = 2000
)

func TestNFTDeployment(t *testing.T) {
	b := newEmulator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode := loadNonFungibleToken()
	nftAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

	// Should be able to deploy a contract as a new account with no keys.
	tokenCode := loadKittyItems(nftAddr.String())
	_, err = b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "KittyItems",
			Source: string(tokenCode),
		},
	})
	if !assert.NoError(t, err) {
		t.Log(err.Error())
	}
	_, err = b.CommitBlock()
	assert.NoError(t, err)

}

func TestCreateKittyItem(t *testing.T) {
	b := newEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode := loadNonFungibleToken()
	nftAddr, _ := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})

	// First, deploy the contract
	tokenCode := loadKittyItems(nftAddr.String())
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	tokenAddr, _ := b.CreateAccount([]*flow.AccountKey{tokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "KittyItems",
			Source: string(tokenCode),
		},
	})

	supply := executeScriptAndCheck(t, b, generateInspectKittyItemSupplyScript(nftAddr.String(), tokenAddr.String()), nil)
	assert.Equal(t, cadence.NewUInt64(0), supply.(cadence.UInt64))

	len := executeScriptAndCheck(
		t,
		b,
		generateInspectCollectionLenScript(nftAddr.String(), tokenAddr.String()),
		[][]byte{jsoncdc.MustEncode(cadence.NewAddress(tokenAddr))},
	)
	assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))

	t.Run("Should be able to mint a token", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(generateMintKittyItemScript(nftAddr.String(), tokenAddr.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(tokenAddr)
		tx.AddArgument(cadence.NewAddress(tokenAddr))
		tx.AddArgument(cadence.NewUInt64(typeID1))

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr},
			[]crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// Assert that the account's collection is correct
		len := executeScriptAndCheck(
			t,
			b,
			generateInspectCollectionLenScript(nftAddr.String(), tokenAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(tokenAddr))},
		)
		assert.Equal(t, cadence.NewInt(1), len.(cadence.Int))

		// Assert that the token type is correct
		/*typeID := executeScriptAndCheck(
			t,
			b,
			generateInspectKittyItemTypeIDScript(nftAddr.String(), tokenAddr.String()),
			// Cheat: We know it's token ID 0
			[][]byte{jsoncdc.MustEncode(cadence.NewUInt64(0))},
		)
		assert.Equal(t, cadence.NewUInt64(typeID1), typeID.(cadence.UInt64))*/
	})

	/*t.Run("Shouldn't be able to borrow a reference to an NFT that doesn't exist", func(t *testing.T) {
		// Assert that the account's collection is correct
		result, err := b.ExecuteScript(generateInspectCollectionScript(nftAddr, tokenAddr, tokenAddr, "KittyItems", "KittyItemsCollection", 5), nil)
		require.NoError(t, err)
		assert.True(t, result.Reverted())
	})*/
}

func TestTransferNFT(t *testing.T) {
	b := newEmulator()

	accountKeys := test.AccountKeyGenerator()

	// Should be able to deploy a contract as a new account with no keys.
	nftCode := contracts.NonFungibleToken()
	nftAddr, err := b.CreateAccount(nil, []sdktemplates.Contract{
		{
			Name:   "NonFungibleToken",
			Source: string(nftCode),
		},
	})
	assert.NoError(t, err)

	// First, deploy the contract
	tokenCode := loadKittyItems(nftAddr.String())
	tokenAccountKey, tokenSigner := accountKeys.NewWithSigner()
	tokenAddr, err := b.CreateAccount([]*flow.AccountKey{tokenAccountKey}, []sdktemplates.Contract{
		{
			Name:   "KittyItems",
			Source: string(tokenCode),
		},
	})
	assert.NoError(t, err)

	joshAccountKey, joshSigner := accountKeys.NewWithSigner()
	joshAddress, err := b.CreateAccount([]*flow.AccountKey{joshAccountKey}, nil)

	tx := flow.NewTransaction().
		SetScript(generateMintKittyItemScript(nftAddr.String(), tokenAddr.String())).
		SetGasLimit(100).
		SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
		SetPayer(b.ServiceKey().Address).
		AddAuthorizer(tokenAddr)
	tx.AddArgument(cadence.NewAddress(tokenAddr))
	tx.AddArgument(cadence.NewUInt64(typeID1))

	signAndSubmit(
		t, b, tx,
		[]flow.Address{b.ServiceKey().Address, tokenAddr},
		[]crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
		false,
	)

	// create a new Collection
	t.Run("Should be able to create a new empty NFT Collection", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(generateSetupAccountScript(nftAddr.String(), tokenAddr.String())).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		len := executeScriptAndCheck(
			t,
			b, generateInspectCollectionLenScript(nftAddr.String(), tokenAddr.String()),
			[][]byte{jsoncdc.MustEncode(cadence.NewAddress(joshAddress))},
		)
		assert.Equal(t, cadence.NewInt(0), len.(cadence.Int))

	})

	/*t.Run("Shouldn't be able to withdraw an NFT that doesn't exist in a collection", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(generateTransferScript(nftAddr, tokenAddr, "KittyItems", "KittyItemsCollection", joshAddress, 3)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(tokenAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr},
			[]crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			true,
		)

		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, joshAddress, "KittyItems", "KittyItemsCollection", 0))

		// Assert that the account's collection is correct
		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, tokenAddr, "KittyItems", "KittyItemsCollection", 1))

	})*/

	// transfer an NFT
	/*t.Run("Should be able to withdraw an NFT and deposit to another accounts collection", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(generateTransferScript(nftAddr, tokenAddr, "KittyItems", "KittyItemsCollection", joshAddress, 0)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(tokenAddr)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, tokenAddr},
			[]crypto.Signer{b.ServiceKey().Signer(), tokenSigner},
			false,
		)

		// Assert that the account's collection is correct
		executeScriptAndCheck(t, b, generateInspectCollectionScript(nftAddr, tokenAddr, joshAddress, "KittyItems", "KittyItemsCollection", 0))

		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, joshAddress, "KittyItems", "KittyItemsCollection", 1))

		// Assert that the account's collection is correct
		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, tokenAddr, "KittyItems", "KittyItemsCollection", 0))

	})*/

	// transfer an NFT
	/*t.Run("Should be able to withdraw an NFT and destroy it, not reducing the supply", func(t *testing.T) {
		tx := flow.NewTransaction().
			SetScript(generateDestroyScript(nftAddr, tokenAddr, "KittyItems", "KittyItemsCollection", 0)).
			SetGasLimit(100).
			SetProposalKey(b.ServiceKey().Address, b.ServiceKey().Index, b.ServiceKey().SequenceNumber).
			SetPayer(b.ServiceKey().Address).
			AddAuthorizer(joshAddress)

		signAndSubmit(
			t, b, tx,
			[]flow.Address{b.ServiceKey().Address, joshAddress},
			[]crypto.Signer{b.ServiceKey().Signer(), joshSigner},
			false,
		)

		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, joshAddress, "KittyItems", "KittyItemsCollection", 0))

		// Assert that the account's collection is correct
		executeScriptAndCheck(t, b, generateInspectCollectionLenScript(nftAddr, tokenAddr, tokenAddr, "KittyItems", "KittyItemsCollection", 0))

		executeScriptAndCheck(t, b, generateInspectNFTSupplyScript(nftAddr, tokenAddr, "KittyItems", 1))

	})*/
}

func replaceKittyItemsAddressPlaceholders(code, nftAddress, kittyItemsAddress string) []byte {
	return []byte(replaceStrings(
		code,
		map[string]string{
			nftAddressPlaceholder:        "0x" + nftAddress,
			kittyItemsAddressPlaceHolder: "0x" + kittyItemsAddress,
		},
	))
}

func loadNonFungibleToken() []byte {
	return nft_contracts.NonFungibleToken()
}

func loadKittyItems(nftAddr string) []byte {
	return []byte(strings.ReplaceAll(
		string(readFile(kittyItemsKittyItemsPath)),
		nftAddressPlaceholder,
		"0x"+nftAddr,
	))
}

func generateSetupAccountScript(nftAddr, kittyItemsAddr string) []byte {
	return replaceKittyItemsAddressPlaceholders(
		string(readFile(kittyItemsSetupAccountPath)),
		nftAddr,
		kittyItemsAddr,
	)
}

func generateMintKittyItemScript(nftAddr, kittyItemsAddr string) []byte {
	return replaceKittyItemsAddressPlaceholders(
		string(readFile(kittyItemsMintKittyItemPath)),
		nftAddr,
		kittyItemsAddr,
	)
}

func generateInspectKittyItemSupplyScript(nftAddr, kittyItemsAddr string) []byte {
	return replaceKittyItemsAddressPlaceholders(
		string(readFile(kittyItemsInspectKittyItemSupplyPath)),
		nftAddr,
		kittyItemsAddr,
	)
}

func generateInspectCollectionLenScript(nftAddr, kittyItemsAddr string) []byte {
	return replaceKittyItemsAddressPlaceholders(
		string(readFile(kittyItemsInspectCollectionLenPath)),
		nftAddr,
		kittyItemsAddr,
	)
}

func generateInspectCollectionIdsScript(nftAddr, kittyItemsAddr string) []byte {
	return replaceKittyItemsAddressPlaceholders(
		string(readFile(kittyItemsInspectCollectionIdsPath)),
		nftAddr,
		kittyItemsAddr,
	)
}
