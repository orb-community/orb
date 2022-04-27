import { ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute } from '@angular/router';
import { Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { PolicyDetailsComponent } from 'app/shared/components/orb/policy/policy-details/policy-details.component';
import { PolicyInterfaceComponent } from 'app/shared/components/orb/policy/policy-interface/policy-interface.component';

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
    interface: false,
  };

  @ViewChild(PolicyDetailsComponent)
  detailsComponent: PolicyDetailsComponent;

  @ViewChild(PolicyInterfaceComponent)
  interfaceComponent: PolicyInterfaceComponent;

  constructor(
    private route: ActivatedRoute,
    private policiesService: AgentPoliciesService,
    private cdr: ChangeDetectorRef,
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
    this.editMode.interface = false;
  }

  save() {
    // get values from all modified sections' forms and submit through service.
    this.discard();
    this.cdr.markForCheck();
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
