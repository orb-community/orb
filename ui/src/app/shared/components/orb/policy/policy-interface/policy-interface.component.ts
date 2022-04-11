import { Component, Input, OnInit } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-policy-interface',
  templateUrl: './policy-interface.component.html',
  styleUrls: ['./policy-interface.component.scss'],
})
export class PolicyInterfaceComponent implements OnInit {
  @Input()
  policy: AgentPolicy = {};

  constructor() { }

  ngOnInit(): void {
  }

}
