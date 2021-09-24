import { Component, OnInit } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.view.component.html',
  styleUrls: ['./agent.view.component.scss'],
})
export class AgentViewComponent implements OnInit {
  strings = STRINGS.agents;

  isLoading: boolean = true;

  agent: Agent;
  agentID;

  command2copy: string;

  constructor(
    private agentsService: AgentsService,
    protected route: ActivatedRoute,
    protected router: Router,
  ) { }
  ngOnInit(): void {
    this.getAgent();
  }

  getAgent() {
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.command2copy = '';

    this.agent = { orb_tags: {}, agent_tags: {} };

    !!this.agentID && this.agentsService.getAgentById(this.agentID).subscribe(resp => {
      this.agent = resp;

      this.command2copy = `docker run --rm --net=host\\
      \r-e ORB_CLOUD_ADDRESS=${document.location.protocol}://${document.location.hostname}/\\
      \r-e ORB_CLOUD_MQTT_ID='${this.agent.id}'\\
      \r-e ORB_CLOUD_MQTT_CHANNEL_ID='${this.agent.channel_id}'\\
      \r-e ORB_BACKENDS_PKTVISOR_IFACE=[ETH-INTERFACE]\\
      \r-e ORB_CLOUD_MQTT_KEY=[AGENT-KEY]\\
      \rns1labs/orb-agent`;
    });
  }
}
