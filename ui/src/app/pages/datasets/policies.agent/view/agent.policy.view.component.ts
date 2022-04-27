import { Component, OnDestroy, OnInit } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

@Component({
  selector: 'ngx-agent-view',
  templateUrl: './agent.policy.view.component.html',
  styleUrls: ['./agent.policy.view.component.scss'],
})
export class AgentPolicyViewComponent implements OnInit, OnDestroy {
  strings = STRINGS.agents;

  isLoading: boolean;

  policyId: string;

  policy: AgentPolicy;

  policySubscription: Subscription;

  editMode = {
    details: false,
    groups: false,
    interface: false,
  };

  constructor(
    protected route: ActivatedRoute,
    protected router: Router,
    protected policiesService: AgentPoliciesService,
  ) {
  }

  ngOnInit() {
    this.policyId = this.route.snapshot.paramMap.get('id');
    this.retrievePolicy();
  }

  isEditMode() {
    return Object.values(this.editMode).reduce((prev, cur) => prev || cur, false);
  }

  discard() {
    this.editMode.details = false;
    this.editMode.groups = false;
    this.editMode.interface = false;
  }

  save() {
    // get values from all modified sections' forms and submit through service.
  }

  retrievePolicy() {
    this.isLoading = true;

    this.policySubscription = this.policiesService
      .getAgentPolicyById(this.policyId)
      .subscribe(policy => {
        this.policy = policy;

        this.isLoading = false;
      });
  }

  ngOnDestroy() {
    this.policySubscription?.unsubscribe();
  }
}
