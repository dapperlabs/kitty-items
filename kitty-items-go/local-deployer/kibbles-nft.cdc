import FungibleToken from 0x9a0766d93b6608b7

pub contract DietKibbles: FungibleToken {
  pub var totalSupply: UFix64

  pub let privatePath: Path
  pub let publicPath: Path

  pub event TokensInitialized(initialSupply: UFix64)
  pub event TokensMinted(amount: UFix64)
  pub event TokensWithdrawn(amount: UFix64, from: Address?)
  pub event TokensDeposited(amount: UFix64, to: Address?)

  pub resource Vault: FungibleToken.Provider, FungibleToken.Receiver, FungibleToken.Balance {
    pub var balance: UFix64

    init(balance: UFix64) {
      self.balance = balance
    }

    pub fun withdraw(amount: UFix64): @FungibleToken.Vault {
      self.balance = self.balance - amount
      emit TokensWithdrawn(amount: amount, from: self.owner?.address)
      return <-create Vault(balance: amount)
    }

    pub fun deposit(from: @FungibleToken.Vault) {
      let vault <- from as! @DietKibbles.Vault
      self.balance = self.balance + vault.balance
      emit TokensDeposited(amount: vault.balance, to: self.owner?.address)
      vault.balance = 0.0
      destroy vault
    }

    destroy() {
      DietKibbles.totalSupply = DietKibbles.totalSupply - self.balance
    }
  }

  pub fun createEmptyVault(): @FungibleToken.Vault {
    return <-create Vault(balance: 0.0)
  }

  pub fun hasDietKibbles(_ address: Address): Bool {
    return getAccount(address)
      .getCapability<&{FungibleToken.Receiver, FungibleToken.Balance}>(DietKibbles.publicPath)!
      .check()
  }

  pub fun fetchBalance(_ address: Address): UFix64? {
    let cap = getAccount(address)
      .getCapability<&{FungibleToken.Balance}>(DietKibbles.publicPath)!

    if let bbb = cap.borrow() {
      return bbb.balance
    }

    return nil
  }

  // Temporary...
  pub fun mintTenDietKibbles(): @DietKibbles.Vault {
    let amount: UFix64 = 10.0
    DietKibbles.totalSupply = DietKibbles.totalSupply + amount
    emit TokensMinted(amount: amount)
    return <-create Vault(balance: amount)
  }

  init() {
    self.totalSupply = 0.0
    self.privatePath = /storage/dietKibbles
    self.publicPath = /public/dietKibbles

    emit TokensInitialized(initialSupply: self.totalSupply)
  }
}