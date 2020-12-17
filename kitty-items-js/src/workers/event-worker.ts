import { parentPort, workerData } from 'worker_threads';
import {config} from "@onflow/config"
import {getEvents} from "@onflow/sdk-build-get-events"
import {send} from "@onflow/sdk-send"
import {decode} from "@onflow/sdk-decode"

config()
  .put("accessNode.api", workerData.flowNode);

parentPort?.on('message', (event) => {
  Promise.resolve(
    send([
      getEvents(workerData.eventType, event.fromBlock, event.toBlock)
    ])
  )
  .then(decode)
  .then(data => {
    if (data && data.length > 0) {
      parentPort?.postMessage({
        data
      });
    }
  })
  .catch(e => {
    console.error(e)
  });
});
