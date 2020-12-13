import * as dotenv from "dotenv";
import * as fcl from "@onflow/fcl";
import initApp from "./app";
import { KibblesService } from "./services/kibbles";
import { FlowService } from "./services/flow";
import {
  fungibleTokenContractAddressFor,
  parseFlowEnvironment,
} from "./constants";

async function run() {
  dotenv.config();
  fcl.config().put("accessNode.api", process.env.FLOW_NODE);
  const flowService = new FlowService(
    process.env.MINTER_FLOW_ADDRESS!,
    process.env.MINTER_PRIVATE_KEY!,
    process.env.MINTER_ACCOUNT_KEY_IDX!
  );

  const flowEnv = parseFlowEnvironment(process.env.FLOW_ENV!);
  const fungibleTokenContractAddress = fungibleTokenContractAddressFor(flowEnv);

  const kibblesService = new KibblesService(flowService, {
    fungibleTokenAddress: fungibleTokenContractAddress,
    kibbleContractAddress: process.env.KIBBLE_CONTRACT_ADDRESS!,
  });

  const app = initApp(kibblesService);

  app.listen(3000, () => {
    console.log("Listening on port 3000!");
  });
}

run().catch((e) => {
  console.error("error", e);
  process.exit(1);
});
