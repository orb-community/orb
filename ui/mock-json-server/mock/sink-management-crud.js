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


const createSinkManagement = (name = null, config = {
        description: null,
        tags: null,
        status: null,
        error: null,
        backend: null,
        config: null, // {remote_host: null, username: null}
        ts_created: null
    }) => {
        return {
            id: uuid.v4(),
            name: name ? name : Faker.company.companyName(),
            description: config.description ? config.description : Faker.company.bs(),
            tags: config.tags ? config.tags : [Faker.hacker.adjective(), Faker.hacker.noun()],
            backend: config.type ? config.type : Faker.hacker.ingverb(),
            status: config.status ? config.status : ['active', 'error'][Math.floor(Math.random() * 100) % 2],
            config: config.config ? config.config : {
                remote_host: Faker.internet.domainName(),
                username: Faker.internet.userName(),
            },
            ts_created: createTimeStamp(),
        }
    }
;

module.exports = {
    getSinkManagementList,
    getSinkManagementById,
    setSinkManagementList,
    updateOrCreateSinkManagementItem,
    deleteSinkManagementItem,
    createSinkManagement,
};
