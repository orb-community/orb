import { Component, Input, OnInit } from '@angular/core';
import { AvailableOS } from 'app/common/services/agents/agents.service';
import { Agent, AgentStates } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-agent-provisioning',
  templateUrl: './agent-provisioning.component.html',
  styleUrls: ['./agent-provisioning.component.scss'],
})
export class AgentProvisioningComponent implements OnInit {
  @Input() agent: Agent;

  agentStates = AgentStates;

  command2copy: string;

  copyCommandIcon: string;

  availableOS = [AvailableOS.DOCKER];

  selectedOS = AvailableOS.DOCKER;

  command2show: string;

  hideCommand: boolean;

  constructor() {
    this.command2copy = '';
    this.command2show = '';
    this.copyCommandIcon = 'clipboard-outline';
  }

  ngOnInit(): void {
    this.hideCommand = this.agent?.state !== this.agentStates.new;
    this.makeCommand2Copy();
  }

  toggleIcon(target) {
    if (target === 'command') {
      this.copyCommandIcon = 'checkmark-outline';
    }
  }

  makeCommand2Copy() {
    // TODO: future - store this elsewhere
    if (this.selectedOS === AvailableOS.DOCKER) {
      this.command2copy = `docker run -d --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent?.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent?.channel_id } \\
-e ORB_CLOUD_MQTT_KEY="AGENT_KEY" \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=mock \\
ns1labs/orb-agent:develop`;

      this.command2show = `docker run -d --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent?.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent?.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=<mark>AGENT_KEY</mark> \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>mock</mark> \\
ns1labs/orb-agent:develop`;
    }
  }

  toggleProvisioningCommand() {
    this.hideCommand = !this.hideCommand;
  }
}
