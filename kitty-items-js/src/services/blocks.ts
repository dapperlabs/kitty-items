import path from 'path';
import { Worker } from 'worker_threads';
import { EventsService } from './events';

class BlocksService {
  
  constructor(
    private readonly eventsService: EventsService,
    private readonly flowNode: string,
  ) {}

  initWorker() {
    let blockWorker = new Worker(path.resolve(__dirname, '../workers/worker.js'), {
      workerData: {
        path: './block-worker.ts',
        flowNode: this.flowNode
      }
    });

    blockWorker.on('message', (step) => {
      console.log(step);
      this.eventsService.fetchEvents(step.fromBlock, step.toBlock);
    });
  }
}

export { BlocksService };