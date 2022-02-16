import { Component, Input, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Agent } from 'app/common/interfaces/orb/agent.interface';

@Component({
  selector: 'ngx-agent-key-component',
  templateUrl: './agent.key.component.html',
  styleUrls: ['./agent.key.component.scss'],
})
export class AgentKeyComponent implements OnInit {
  strings = STRINGS.agents;

  command2copy: string;

  command2show: string;
  copyCommandIcon: string;

  key2copy: string;
  copyKeyIcon: string;

  @Input() agent: Agent = {};

  constructor(
    protected dialogRef: NbDialogRef<AgentKeyComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }

  ngOnInit(): void {
    this.makeCommand2Copy();
    this.key2copy = this.agent.key;
    this.copyCommandIcon = 'clipboard-outline';
    this.copyKeyIcon = 'clipboard-outline';
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

  toggleIcon (target) {
    if (target === 'key') {
      this.copyKeyIcon = 'checkmark-outline';
    } else if (target === 'command') {
      this.copyCommandIcon = 'checkmark-outline';
    }
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
