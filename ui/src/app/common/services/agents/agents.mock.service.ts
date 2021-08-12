import {HttpClient, HttpParams} from '@angular/common/http';
import {Injectable} from '@angular/core';
import {Observable} from 'rxjs';
import 'rxjs/add/observable/empty';
import 'rxjs/add/observable/of';
import {Router} from '@angular/router';

import {NotificationsService} from '../../../common/services/notifications/notifications.service';
import {Agent} from 'app/common/interfaces/orb/agent.interface';
import {PageFilters} from 'app/common/interfaces/mainflux.interface';
import * as uuid from 'uuid';


/*TODO
 Mocked service, to be replace by actual service when the BE is ready;
// **/

// default filters
const defLimit: number = 20;


let agentList: Agent[] = [];

const getAgentsList = () => agentList;

const getAgentById = id => agentList.find(
  elem => elem.id === id);

const updateAgentItem = (agentItem) => {
  const index = agentList.findIndex(entry => entry.id === agentItem.id);
  agentList[index] = agentItem;
  agentList = Array.from(agentList);
  return agentList;
};

const deleteAgentItem = (sinkItem) => {
  const index = agentList.findIndex(entry => entry.id === sinkItem.id);
  if (index === -1) {
    return;
  }
  agentList.splice(index, 1);
  agentList = Array.from(agentList);
  return agentList;
};

const addAgentItem = (agentItem: Agent) => {
  agentItem.id = uuid.v4();
  agentList.push(agentItem);
};

@Injectable()
export class AgentsMockService {
  picture = 'assets/images/mainflux-logo.png';

  constructor(
    private router: Router,
    private notificationsService: NotificationsService,
  ) {
  }

  addAgent(agentItem: Agent) {
    addAgentItem(agentItem);
    return Observable.of(agentItem);
  }

  getAgentById(agentId: string): any {
    const agent = getAgentById(agentId);
    return Observable.of(agent);
  }

  getAgents(filters: PageFilters) {
    const list = getAgentsList();
    const reply = {
        agents: list.slice(filters.offset, filters.limit),
        total: list.length,
        offset: filters.offset || 0,
        limit: filters.limit || defLimit,
    };

    return Observable.of(reply);
  }

  editSink(agentItem: Agent): any {
    const agent = updateAgentItem(agentItem);
    return Observable.of(agent);
  }

  deleteAgent(agentId: string) {
    const agent = getAgentById(agentId);
    deleteAgentItem(agentId);
    return Observable.of(agent);
  }
}
