const path = require('path');
const { workerData } = require('worker_threads');

try {
  require('ts-node').register();
  require(path.resolve(__dirname, workerData.path));
} catch(error) {
  console.error(error)
}