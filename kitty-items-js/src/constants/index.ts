enum FlowEnvironment {
  Emulator,
  Testnet,
  Mainnet,
}

export function fungibleTokenContractAddressFor(
  environment: FlowEnvironment
): string {
  switch (environment) {
    case FlowEnvironment.Emulator:
      return "0xee82856bf20e2aa6";
    case FlowEnvironment.Testnet:
      return "0x9a0766d93b6608b7";
    case FlowEnvironment.Mainnet:
      return "0xf233dcee88fe0abe";
  }
  return "";
}

export function parseFlowEnvironment(environment: string) {
  switch (environment) {
    case "emulator":
      return FlowEnvironment.Emulator;
    case "testnet":
      return FlowEnvironment.Testnet;
    case "mainnet":
      return FlowEnvironment.Mainnet;
    default:
      throw new Error(`invalid flow environment: ${environment}`);
  }
}
