import { parentPort, workerData } from 'worker_threads';
import {config} from "@onflow/config"
import {getEvents} from "@onflow/sdk-build-get-events"
import {send} from "@onflow/sdk-send"
import {decode} from "@onflow/sdk-decode"

config()
  .put("accessNode.api", workerData.flowNode);

parentPort?.on('message', (step) => {
  Promise.resolve(
    send([
      getEvents(workerData.eventType, step.fromBlock, step.toBlock)
    ])
  )
  .then(decode)
  .then(events => {
    if (events && events.length > 0) {
      parentPort?.postMessage(events);
    }
  })
  .catch(e => {
    console.error(e)
  });
});
