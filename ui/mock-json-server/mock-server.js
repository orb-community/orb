// This is responsible to create the mock server.
const express = require('express');
const bodyParser = require('body-parser');
const cors = require('cors');
const fs = require('fs');

const SINK_TMP_FILE = `${__dirname}/mock/sinks.json`;
const {
  getSinkManagementList,
  getSinkManagementById,
  setSinkManagementList,
  updateOrCreateSinkManagementItem,
  deleteSinkManagementItem,
  createSinkManagement,
} = require('./mock/sinks-crud');

process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';

const writeTmpFile = (items, filePath) => {
  fs.writeFileSync(filePath, JSON.stringify(items));
};

const getLatestTmp = (listSetter, filePath) => {
  if (fs.existsSync(filePath)) {
    const json = fs.readFileSync(filePath);
    listSetter(JSON.parse(json));
  } else {
    writeTmpFile([], filePath);
  }
};


const app = express();
const PORT = 3000;

const MOCK_DELAY = 1000;

function send(callback, delay = MOCK_DELAY) {
  if (MOCK_DELAY) {
    setTimeout(callback, delay);
  } else {
    callback();
  }
}

app.use(
  bodyParser.urlencoded({
    extended: true,
  }),
);
app.use(bodyParser.json());
app.use(cors());

app.use((req, res, next) => {
  res.header('Access-Control-Allow-Origin', '*');
  res.header(
    'Access-Control-Allow-Headers',
    'Origin, X-Requested-With, Content-Type, Accept',
  );
  console.log(`${req.method} ${req.url}`);
  next();
});

// Sink Management Crud

app.get('/sinks', (req, res) => {
  console.log('gotten list of sinks');
  send(() => res.send(getSinkManagementList()), MOCK_DELAY);
});

app.get('/sinks/:id', (req, res) => {
  const { id } = req.params;
  console.log('gotten a single sink');
  const sinkItem = getSinkManagementById(id);
  send(() => res.send(sinkItem), MOCK_DELAY);
});

app.post('/sinks', (req, res) => {
  const data = req.body;
  const newSinkItem = createSinkManagement(data.name);
  const newList = updateOrCreateSinkManagementItem(newSinkItem);
  setSinkManagementList(newList);
  writeTmpFile(newList, SINK_TMP_FILE);
   console.log('created a sink');
  send(() => res.send(newSinkItem), MOCK_DELAY);
});

app.put('/sinks/:id', (req, res) => {
  const { id } = req.params;
  const sinkPayload = req.body;
  const sinkItem = getSinkManagementById(id);
  if (!sinkItem) {
    throw new Error(`Error while updating sink with id ${id}`);
  }
  const newList =  updateOrCreateSinkManagementItem(sinkPayload);
  writeTmpFile(newList, SINK_TMP_FILE);
  console.log('updated a sink');
  send(() => res.sendStatus(204), MOCK_DELAY);
});

app.delete('/sinks/:id', (req, res) => {
  const {id} = req.params;
  let sinkItem = getSinkManagementById(id);
  if (!sinkItem) {
    throw new Error(`Error deleting sink with id ${id}`);
  }

  const newList = deleteSinkManagementItem(sinkItem);
  console.log('deleted a sink');
  writeTmpFile(newList, SINK_TMP_FILE);
});

const key = fs.readFileSync(`${__dirname}/certs/selfsigned.key`);
const cert = fs.readFileSync(`${__dirname}/certs/selfsigned.crt`);
const options = {
  key,
  cert,
};

getLatestTmp(setSinkManagementList, SINK_TMP_FILE);

app.listen(PORT, () => {
  console.log('JSON Server is running on port:', PORT);
});
