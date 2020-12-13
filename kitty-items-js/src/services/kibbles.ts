import * as t from "@onflow/types";
import * as fcl from "@onflow/fcl";
import { FlowService } from "./flow";
import { mintKibblesTemplate } from "../templates";

interface KibblesMinterParams {
  fungibleTokenAddress: string;
  kibbleContractAddress: string;
}

class KibblesService {
  constructor(
    private readonly flowService: FlowService,
    private readonly kibblesParams: KibblesMinterParams
  ) {}
  async mintKibblesToAddress(
    destinationAddress: string,
    amount: string
  ): Promise<string> {
    const authorization = this.flowService.authorizeMinter();

    const template = mintKibblesTemplate
      .replace(
        /0xFUNGIBLETOKENADDRESS/gi,
        `0x${this.kibblesParams.fungibleTokenAddress}`
      )
      .replace(/0xKIBBLE/gi, `0x${this.kibblesParams.kibbleContractAddress}`);

    const response = await fcl.send([
      fcl.transaction`
        ${template}
      `,
      fcl.args([
        fcl.arg(destinationAddress, t.Address),
        fcl.arg(amount, t.UFix64),
      ]),
      fcl.proposer(authorization),
      fcl.payer(authorization),
      fcl.authorizations([authorization]),
      fcl.limit(100),
    ]);

    return await fcl.tx(response).onceExecuted();
  }
}

export { KibblesService };
