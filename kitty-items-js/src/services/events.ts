import path from 'path';
import { Worker } from 'worker_threads';

interface Event {
  type: string,
  transactionId: string,
  transactionIndex: number,
  eventIndex: number,
  data:any
}

interface EventHandler {
  key: string
  callback: (event: Event) => void
}

class EventsService {

  workers: Worker[];

  kibblesEvents: EventHandler[] = [
    {
      key: "TokensWithdrawn",
      callback: this.kibbleTokensWithdrawn
    },
    {
      key: "TokensDeposited",
      callback: this.kibbleTokensDeposited
    },
    {
      key: "TokensMinted",
      callback: this.kibbleTokensMinted
    },
    {
      key: "TokensBurned",
      callback: this.kibbleTokensBurned
    },
    {
      key: "MinterCreated",
      callback: this.kibbleMinterCreated
    },
    {
      key: "BurnerCreated",
      callback: this.kibbleBurnerCreated
    }
  ];

  kittyItemsEvents: EventHandler[] = [
    {
      key: "Withdraw",
      callback: this.kittyItemWithdrawn
    },
    {
      key: "Deposit",
      callback: this.kittyItemDeposited
    },
  ];

  marketEvents: EventHandler[] = [
    {
      key: "SaleOfferCreated",
      callback: this.saleOfferCreated
    },
    {
      key: "SaleOfferAccepted",
      callback: this.saleOfferAccepted
    },
    {
      key: "SaleOfferFinished",
      callback: this.saleOfferFinished
    },
    {
      key: "CollectionInsertedSaleOffer",
      callback: this.collectionInsertedSaleOffer
    },
    {
      key: "CollectionRemovedSaleOffer",
      callback: this.collectionRemovedSaleOffer
    },
  ]

  constructor(
    private readonly contractAddress: string,
    private readonly flowNode: string,
  ) {
    this.workers = [];
    this.initWorkers();
  }

  initWorkers = () => {
    this.initContractWorker('Kibble', this.kibblesEvents);
    this.initContractWorker('KittyItems', this.kittyItemsEvents);
    this.initContractWorker('KittyItemsMarket', this.marketEvents);
  }

  initContractWorker = (contractName: string, eventsList: EventHandler[]) => {
    for (let i = 0; i < eventsList.length; i++) {
      let worker = new Worker(path.resolve(__dirname, '../workers/worker.js'), {
        workerData: {
          path: './event-worker.ts',
          flowNode: this.flowNode,
          eventType: `A.${this.contractAddress}.${contractName}.${eventsList[i].key}`
        }
      }); 
      worker.on('message', async (events) => {
        for (let event of events) {
          await eventsList[i].callback(event)
        }
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

  async kibbleTokensWithdrawn(event) {
    console.log("kibbleTokensWithdrawn", event)
  }

  async kibbleTokensDeposited(event) {
    console.log("kibbleTokensDeposited", event)
  }

  async kibbleTokensMinted(event) {
    console.log("kibbleTokensMinted", event)
  }

  async kibbleTokensBurned(event) {
    console.log("kibbleTokensBurned", event)
  }

  async kibbleMinterCreated(event) {
    console.log("kibbleMinterCreated", event)
  }

  async kibbleBurnerCreated(event) {
    console.log("kibbleBurnerCreated", event)
  }

  async kittyItemWithdrawn(event) {
    console.log("kittyItemWithdrawn", event)
  }

  async kittyItemDeposited(event) {
    console.log("kittyItemDeposited", event)
  }

  async saleOfferCreated(event) {
    console.log("saleOfferCreated", event)
  }

  async saleOfferAccepted(event) {
    console.log("saleOfferAccepted", event)
  }

  async saleOfferFinished(event) {
    console.log("saleOfferFinished", event)
  }

  async collectionInsertedSaleOffer(event) {
    console.log("collectionInsertedSaleOffer", event)
  }

  async collectionRemovedSaleOffer(event) {
    console.log("collectionRemovedSaleOffer", event)
  }

}

export { EventsService };