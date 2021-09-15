import { Injectable } from '@angular/core';
import { of } from 'rxjs';
import 'rxjs/add/observable/empty';
import 'rxjs/add/observable/of';

import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface.ts';
import { PageFilters } from 'app/common/interfaces/mainflux.interface';
import * as Moment from 'moment';
import * as Faker from 'faker';
import * as uuid from 'uuid';


/*TODO
 Mocked service, to be replace by actual service when the BE is ready;
// **/

// default filters
const defLimit: number = 20;


let agentGroupList: AgentGroup[] = [];
// const agentList: Agent[] = [];

const getTimeStamp = () => Moment().toISOString();

const genAgent = () => ({
  id: uuid.v4(),
  key: uuid.v4(),
  channel_id: uuid.v4(),
  name: Faker.company.companyName(),
  state: 'active',
});

const getAgentList = (size) => {
  const list = [];
  for (let i = 0; i < size; i++) {
    list.push(genAgent());
  }
  return list;
};

const genGroupList = () => {
  const agents = getAgentList(Faker.datatype.number({ min: 3, max: 9 }));
  return ({
    id: uuid.v4(),
    name: Faker.company.companyName(),
    description: Faker.lorem.words(),
    tags: {
      cloud: 'aws',
    },
    ts_created: getTimeStamp(),
    matching_agents: {
      total: agents.length,
      online: Faker.datatype.number({ min: agents.length / 3, max: agents.length }),
    },
    agents: agents,
  });
};

const getAgentsGroupList = () => agentGroupList;

const getAgentGroupById = id => agentGroupList.find(
  elem => elem.id === id);

const updateAgentGroupItem = (agentItem) => {
  const index = agentGroupList.findIndex(entry => entry.id === agentItem.id);
  agentGroupList[index] = agentItem;
  agentGroupList = Array.from(agentGroupList);
  return agentGroupList;
};

const deleteAgentGroupItem = (sinkItem) => {
  const index = agentGroupList.findIndex(entry => entry.id === sinkItem.id);
  if (index === -1) {
    return;
  }
  agentGroupList.splice(index, 1);
  agentGroupList = Array.from(agentGroupList);
  return agentGroupList;
};

const addAgentGroupItem = (agentItem: AgentGroup) => {
  agentItem.id = uuid.v4();
  agentGroupList.push(agentItem);
};

@Injectable()
export class AgentsMockService {

  constructor() {
    for (let i = 0; i < 10; i++) {
      const group = genGroupList();
      agentGroupList.push(group);
    }
  }

  addAgentGroup(agentItem: AgentGroup) {
    addAgentGroupItem(agentItem);

    return of(agentItem);
  }

  getAgentGroupById(agentId: string): any {
    const agent = getAgentGroupById(agentId);
    return of(agent);
  }

  getAgentGroups(filters: PageFilters) {
    const list = getAgentsGroupList();
    const reply = {
      agents: list.slice(filters.offset, filters.limit),
      total: list.length,
      offset: filters.offset || 0,
      limit: filters.limit || defLimit,
    };

    return of(reply);
  }

  checkAgents(args) {
    const newAgentgroup: AgentGroup = {
      id: uuid.v4(),
      name: Faker.company.companyName(),
      description: Faker.lorem.words(),
      tags: [],
      matching_agents: {
        total: 1,
        online: 2,
      },
      agents: [],
    };

    return of(newAgentgroup);
  }

  editAgentGroup(agentItem: Agent): any {
    const agent = updateAgentGroupItem(agentItem);
    return of(agent);
  }

  deleteAgentGroup(agentId: string) {
    const agent = getAgentGroupById(agentId);
    deleteAgentGroupItem(agentId);
    return of(agent);
  }
}
