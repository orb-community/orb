import { Component, Input, OnInit } from '@angular/core';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-agent-backends',
  templateUrl: './agent-backends.component.html',
  styleUrls: ['./agent-backends.component.scss'],
})
export class AgentBackendsComponent implements OnInit {
  @Input() agent: Agent;

  constructor(
    protected notificationService: NotificationsService,
  ) {

  }

  ngOnInit(): void {}

}
