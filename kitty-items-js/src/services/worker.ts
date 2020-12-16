import { parentPort, workerData } from 'worker_threads';
 
setInterval(() => {
  if(parentPort) {
    parentPort.postMessage(
      "Hello"
    );
  }
}, 1000);