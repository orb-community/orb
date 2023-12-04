import { Component, Input, OnInit } from '@angular/core';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-agent-backends',
  templateUrl: './agent-backends.component.html',
  styleUrls: ['./agent-backends.component.scss'],
})
export class AgentBackendsComponent implements OnInit {
  @Input() agent: Agent;

  agentStates = AgentStates;

  identify(index, item) {
    return item.id;
  }

  constructor(
    protected notificationService: NotificationsService,
  ) {

  }

  ngOnInit(): void {}

}
