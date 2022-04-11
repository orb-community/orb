import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';

@Component({
  selector: 'ngx-agent-policy-details-component',
  templateUrl: './agent.policy.details.component.html',
  styleUrls: ['./agent.policy.details.component.scss'],
})
export class AgentPolicyDetailsComponent implements OnInit, OnDestroy {
  @Input() agentPolicy: AgentPolicy = {};

  isLoading: boolean;

  subscription: Subscription;

  constructor(
    protected dialogRef: NbDialogRef<AgentPolicyDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
    protected agentPoliciesService: AgentPoliciesService,
  ) {
    this.isLoading = true;
  }

  ngOnInit() {
    const { id } = this.agentPolicy;
    this.subscription = this.agentPoliciesService
      .getAgentPolicyById(id)
      .subscribe(agentPolicy => {
        this.agentPolicy = agentPolicy;
        this.isLoading = false;
      });
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onOpenEdit() {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
