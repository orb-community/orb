import { NgModule } from '@angular/core';

import 'rxjs/add/operator/catch';
import 'rxjs/add/observable/throw';
import 'rxjs/add/operator/switchMap';
import 'rxjs/add/operator/map';

import { BootstrapService } from './services/bootstrap/bootstrap.service';
import { ThingsService } from './services/things/things.service';
import { UsersService } from './services/users/users.service';
import { UserGroupsService } from './services/users/groups.service';
import { FsService } from './services/fs/fs.service';
import { IntervalService } from './services/interval/interval.service';
// Orb
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { FleetsService } from 'app/common/services/fleets/fleets.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';

import { TokenInterceptor } from 'app/auth/auth.token.interceptor.service';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';

// Orb mock

@NgModule({
  providers: [
    BootstrapService,
    ThingsService,
    NotificationsService,
    UsersService,
    UserGroupsService,
    FsService,
    IntervalService,
    // ORB Services
    AgentPoliciesService,
    AgentGroupsService,
    AgentsService,
    DatasetPoliciesService,
    FleetsService,
    SinksService,
    OrbService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: TokenInterceptor,
      multi: true,
    },
  ],
})
export class CommonModule {}
