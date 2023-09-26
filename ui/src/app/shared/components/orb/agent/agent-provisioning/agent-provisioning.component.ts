import { Component, Input, OnInit } from "@angular/core";
import { AvailableOS } from "app/common/services/agents/agents.service";
import { Agent, AgentStates } from "app/common/interfaces/orb/agent.interface";

@Component({
  selector: "ngx-agent-provisioning",
  templateUrl: "./agent-provisioning.component.html",
  styleUrls: ["./agent-provisioning.component.scss"],
})
export class AgentProvisioningComponent implements OnInit {
  @Input() agent: Agent;

  agentStates = AgentStates;

  copyCommandIcon: string;

  availableOS = [AvailableOS.DOCKER];

  selectedOS = AvailableOS.DOCKER;

  defaultCommandCopy: string;
  defaultCommandShow: string;
  fileConfigCommandCopy: string;
  fileConfigCommandShow: string;

  hideCommand: boolean;
  hideCommand2: boolean;
  hideCommand3: boolean;

  constructor() {
    this.copyCommandIcon = "copy-outline";
  }

  ngOnInit(): void {
    this.hideCommand2 = false;
    this.hideCommand3 = true;
    this.hideCommand = this.agent?.state !== this.agentStates.new;
    this.makeCommand2Copy();
  }

  toggleIcon(target) {
    if (target === "command") {
      this.copyCommandIcon = "checkmark-outline";
      setTimeout(() => {
        this.copyCommandIcon = "copy-outline";
      }, 2000);
    }
  }

  makeCommand2Copy() {
    this.defaultCommandCopy = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY="AGENT_KEY" \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=auto \\
orbcommunity/orb-agent`;

    this.defaultCommandShow = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=<mark>AGENT_KEY</mark> \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>auto</mark> \\
orbcommunity/orb-agent`;

  this.fileConfigCommandCopy = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY="AGENT_KEY" \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=auto \\
-v \${PWD}/:/usr/local/orb/ \\
orbcommunity/orb-agent run -c /usr/local/orb/agent.yaml`;

  this.fileConfigCommandShow = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=<mark>AGENT_KEY</mark> \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>auto</mark> \\
-v \${PWD}/:/usr/local/orb/ \\
orbcommunity/orb-agent run -c /usr/local/orb/agent.yaml`;
  }

toggleProvisioningCommand(command: string) {
  switch (command) {
    case 'hideCommand':
      this.hideCommand = !this.hideCommand;
      break;
    case 'hideCommand2':
      this.hideCommand2 = !this.hideCommand2;
      break;
    case 'hideCommand3':
      this.hideCommand3 = !this.hideCommand3;
      break;
  }
}
}
