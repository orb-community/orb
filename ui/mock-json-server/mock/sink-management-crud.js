const uuid = require('uuid');
const Moment = require('moment');
const Faker = require('faker');

let sinkManagementList = [];

const createTimeStamp = () => Moment().toISOString();

const getSinkManagementList = () => sinkManagementList;

const getSinkManagementById = id => sinkManagementList.find(elem => elem.id === id);

const setSinkManagementList = list => sinkManagementList = list;

const updateOrCreateSinkManagementItem = (sinkItem) => {
  sinkItem.updatedAt = createTimeStamp();
  const index = sinkManagementList.findIndex(entry => entry.id === sinkItem.id);
  if (index === -1) {
    sinkManagementList.push(sinkItem);
    return sinkManagementList;
  }

  sinkManagementList[index] = sinkItem;
  sinkManagementList = Array.from(sinkManagementList);
  return sinkManagementList;
};

const deleteSinkManagementItem = (sinkItem) => {
  const index = sinkManagementList.findIndex(entry => entry.id === sinkItem.id);
  if (index === -1) {
    return;
  }
  sinkManagementList.splice(index, 1);
  sinkManagementList = Array.from(sinkManagementList);
  return sinkManagementList;
}


const createSinkManagement = (name = null, config = null) => {
  return {
    id: uuid.v4(),
    name: name ? name : Faker.company.name(),
    mfOwnerId: uuid.v4(),
    metadata: 'sample',
    createdAt: createTimeStamp(),
  }
};

module.exports = {
  getSinkManagementList,
  getSinkManagementById,
  setSinkManagementList,
  updateOrCreateSinkManagementItem,
  deleteSinkManagementItem,
  createSinkManagement,
};
