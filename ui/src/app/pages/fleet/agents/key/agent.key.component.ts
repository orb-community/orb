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

  defaultCommandCopy: string;
  defaultCommandShow: string;
  fileConfigCommandCopy: string;
  fileConfigCommandShow: string;

  copyCommandIcon: string;
  copyCommandIcon2: string;

  key2copy: string;
  copyKeyIcon: string;
  saveKeyIcon: string;

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
    this.copyCommandIcon = 'copy-outline';
    this.copyCommandIcon2 = 'copy-outline';
    this.copyKeyIcon = 'copy-outline';
    this.saveKeyIcon = 'save-outline';
  }

  makeCommand2Copy() {
    this.defaultCommandCopy = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=auto \\
orbcommunity/orb-agent`;

    this.defaultCommandShow = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \\
-e PKTVISOR_PCAP_IFACE_DEFAULT=<mark>auto</mark> \\
orbcommunity/orb-agent`;

  this.fileConfigCommandCopy = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \\
-v \${PWD}/:/opt/orb/ \\
orbcommunity/orb-agent run -c /opt/orb/agent.yaml`;

  this.fileConfigCommandShow = `docker run -d --restart=always --net=host \\
-e ORB_CLOUD_ADDRESS=${ document.location.hostname } \\
-e ORB_CLOUD_MQTT_ID=${ this.agent.id } \\
-e ORB_CLOUD_MQTT_CHANNEL_ID=${ this.agent.channel_id } \\
-e ORB_CLOUD_MQTT_KEY=${ this.agent.key } \\
-v \${PWD}/:/opt/orb/ \\
orbcommunity/orb-agent run -c /opt/orb/agent.yaml`;
  }

  toggleIcon (target) {
    if (target === 'key') {
      this.copyKeyIcon = 'checkmark-outline';
      setTimeout(() => {
        this.copyKeyIcon = 'copy-outline';
      }, 2000);
    } else if (target === 'command') {
      this.copyCommandIcon = 'checkmark-outline';
      setTimeout(() => {
        this.copyCommandIcon = 'copy-outline';
      }, 2000);
    } else if (target === 'command2') {
      this.copyCommandIcon2 = 'checkmark-outline';
      setTimeout(() => {
        this.copyCommandIcon2 = 'copy-outline';
      }, 2000);
    }
  }

  onClose() {
    this.dialogRef.close(false);
  }
  downloadCommand(commandType: string) {
    if (commandType === 'default') {
      const blob = new Blob([this.defaultCommandCopy], { type: 'text/plain' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${this.agent.id}.txt`;
      a.click();
      window.URL.revokeObjectURL(url);
    } else if (commandType === 'fileConfig') {
      const blob = new Blob([this.fileConfigCommandCopy], { type: 'text/plain' });
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${this.agent.id}_configfile.txt`;
      a.click();
      window.URL.revokeObjectURL(url);
    }

  }

}
