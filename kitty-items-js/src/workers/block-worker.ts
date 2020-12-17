import { parentPort, workerData } from 'worker_threads';
import {config} from "@onflow/config"
import {latestBlock} from "@onflow/sdk-latest-block"

let currentBlockHeight = 18205000;
const stepSize = 1000;
const tickInterval = 5 * 1000;

config()
  .put("accessNode.api", workerData.flowNode)

setInterval(() => {
  Promise.resolve(latestBlock()).then(lastBlock => {
    if (lastBlock.height > currentBlockHeight) {
      const fromBlock = currentBlockHeight;
      let toBlock = currentBlockHeight + stepSize;
      if (toBlock > lastBlock.height) {
        toBlock = lastBlock.height;
      } 
      parentPort?.postMessage({
        fromBlock,
        toBlock
      });
      currentBlockHeight = toBlock
    }  
  })
  .catch(e => {
    console.error(e)
  });
}, tickInterval);
