import { Component } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentsService } from 'app/common/services/agents/agents.service';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.view.component.html',
  styleUrls: ['./agent.view.component.scss'],
})
export class AgentViewComponent {
  strings = STRINGS.agents;

  isLoading: boolean = true;

  agent: Agent;
  agentID;

  command2copy: string;
  copyCommandIcon: string;

  command2show: string;


  constructor(
    private agentsService: AgentsService,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    this.agent = this.router.getCurrentNavigation().extras.state?.agent as Agent || null;
    this.agentID = this.route.snapshot.paramMap.get('id');
    this.command2copy = '';
    this.command2show = '';
    this.copyCommandIcon = 'clipboard-outline';

    !!this.agentID && this.agentsService.getAgentById(this.agentID).subscribe(resp => {
      this.agent = resp;
      this.makeCommand2Copy();
      this.isLoading = false;
    });
  }

  toggleIcon (target) {
  if (target === 'command') {
      this.copyCommandIcon = 'checkmark-outline';
    }
  }

  makeCommand2Copy() {
    this.command2copy = `docker run -d --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=mock \\
ns1labs/orb-agent:develop`;

    this.command2show = `docker run -d --net=host \n
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \n
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \n
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \n
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \n
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>mock</mark> \n

ns1labs/orb-agent:develop`;
  }
}
