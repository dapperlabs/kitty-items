import path from 'path';
import { Worker } from 'worker_threads';
import {latestBlock} from "@onflow/sdk-latest-block"

class EventsService {

  workers: Worker[];

  kibblesEvents = [
    {
      key: "TokensWithdrawn",
      callback: this.processTokensWithdrawn
    },
    {
      key: "TokensDeposited",
      callback: this.processTokensDeposited
    },
    {
      key: "TokensMinted",
      callback: this.processTokensMinted
    },
    {
      key: "TokensBurned",
      callback: this.processTokensBurned
    },
    {
      key: "MinterCreated",
      callback: this.processMinterCreated
    },
    {
      key: "BurnerCreated",
      callback: this.processBurnerCreated
    }
  ];

  constructor(
    private readonly kibbleAddress: string,
    private readonly flowNode: string,
  ) {
    this.workers = [];
    this.initWorkers();
  }

  initWorkers = () => {
    const events = this.kibblesEvents.map(v => v.key);
    for (let i = 0; i < events.length; i++) {
      // workers
      let worker = new Worker(path.resolve(__dirname, '../workers/worker.js'), {
        workerData: {
          path: './event-worker.ts',
          flowNode: this.flowNode,
          eventType: `A.${this.kibbleAddress}.Kibble.${events[i]}`
        }
      }); 
      worker.on('message', (events) => {
        console.log('received events data from event worker',events);
      });
      this.workers.push(worker);
    }
  }

  fetchEvents = (fromBlock, toBlock) => {
    for (let i = 0; i < this.workers.length; i++) {
      this.workers[i].postMessage({
        fromBlock,
        toBlock
      });
    }
  }

  init = async () => {
    let next = await latestBlock()
    console.log(next)
  }

  // Add these to kibble service
  processTokensWithdrawn(event) {
    console.log("TokensWithdrawn")
  }

  processTokensDeposited(event) {
    console.log("processTokensDeposited")
  }

  processTokensMinted(event) {
    console.log("processTokensMinted")
  }

  processTokensBurned(event) {
    console.log("processTokensBurned")
  }

  processMinterCreated(event) {
    console.log("processMinterCreated")
  }

  processBurnerCreated(event) {
    console.log("processBurnerCreated")
  }

}

export { EventsService };