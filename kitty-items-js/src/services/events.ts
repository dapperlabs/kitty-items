import * as fcl from "@onflow/fcl";
import path from 'path';
import { Worker } from 'worker_threads';
import { getEvents } from "@onflow/sdk-build-get-events";

class EventsService {

  worker: Worker;

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

  kibblesEventsList = [
    'TokensWithdrawn',
    'TokensDeposited',
    'TokensMinted',
    'TokensBurned',
    'MinterCreated',
    'BurnerCreated'
  ]

  constructor(
    private readonly kibbleAddress: string
  ) {
    this.worker = new Worker(path.resolve(__dirname, './worker.js'), {
      workerData: {
        path: './worker.ts',
        events: this.kibblesEventsList
      }
    });
    this.worker.on('message', (msg) => {
      console.log(msg);
    })
  }
  
  // subscribeToEvent = (obj) => {
  //   fcl.events(obj.key).subscribe(event => {
  //     console.log("event", event);
  //     // obj.callback(event)
  //   });
    
  // }

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