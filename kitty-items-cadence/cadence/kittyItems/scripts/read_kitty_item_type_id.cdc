import NonFungibleToken from 0xNONFUNGIBLETOKEN
import KittyItems from 0xKITTYITEMS

// This script reads metadata about an NFT in a user's collection
pub fun main(account: Address, itemID: UInt64): UInt64 {

    // Get the public account object of the owner of the token
    let owner = getAccount(account)

    let collectionBorrow = owner
        .getCapability(KittyItems.CollectionPublicPath)!
        .borrow<&KittyItems.Collection{KittyItems.KittyItemsCollectionPublic}>()!

    // Borrow a reference to a specific NFT in the collection
    let kittyItem = collectionBorrow.borrowKittyItem(id: itemID)

    return kittyItem.typeID
}