import {
  ChangeDetectorRef, Component, OnDestroy, OnInit, ViewChild,
} from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import {
  PolicyConfig,
} from 'app/common/interfaces/orb/policy/config/policy.config.interface';
import {
  AgentPoliciesService,
} from 'app/common/services/agents/agent.policies.service';
import {
  NotificationsService,
} from 'app/common/services/notifications/notifications.service';
import {
  PolicyDetailsComponent,
} from 'app/shared/components/orb/policy/policy-details/policy-details.component';
import {
  PolicyInterfaceComponent,
} from 'app/shared/components/orb/policy/policy-interface/policy-interface.component';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';

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
    details: false, interface: false,
  };

  @ViewChild(PolicyDetailsComponent) detailsComponent: PolicyDetailsComponent;

  @ViewChild(
    PolicyInterfaceComponent) interfaceComponent: PolicyInterfaceComponent;

  constructor(
    private route: ActivatedRoute,
    private policiesService: AgentPoliciesService,
    private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
  ) {}

  ngOnInit() {
    this.policyId = this.route.snapshot.paramMap.get('id');
    this.retrievePolicy();
  }

  isEditMode() {
    return Object.values(this.editMode)
      .reduce((prev, cur) => prev || cur, false);
  }

  canSave() {
    const detailsValid = this.editMode.details
      ? this.detailsComponent?.formGroup?.status === 'VALID' : true;

    const interfaceValid = this.editMode.interface
      ? this.interfaceComponent?.formControl?.status === 'VALID' : true;

    return detailsValid && interfaceValid;
  }

  discard() {
    this.editMode.details = false;
    this.editMode.interface = false;
  }

  save() {
    const {
      format, version, name, description, id, tags, backend,
    } = this.policy;

    // get values from all modified sections' forms and submit through service.
    const policyDetails = this.detailsComponent.formGroup?.value;
    const policyInterface = this.interfaceComponent.code;

    // trying to work around rest api
    const detailsPartial = !!this.editMode.details && {
      ...policyDetails,
    } || { name, description };

    let interfacePartial = {};

    if (format === 'yaml') {
      interfacePartial = {
        format,
        policy_data: policyInterface,
      };
    } else {
      interfacePartial = {
        policy: JSON.parse(policyInterface) as PolicyConfig,
      };
    }

    const payload = {
      ...detailsPartial,
      ...interfacePartial,
      version, id, tags, backend,
    } as AgentPolicy;

    this.policiesService.editAgentPolicy(payload)
      .subscribe(resp => {
        this.discard();
        this.retrievePolicy();
        this.cdr.markForCheck();
      });
  }

  retrievePolicy() {
    this.isLoading = true;

    this.policySubscription = this.policiesService
      .getAgentPolicyById(this.policyId)
      .subscribe(policy => {
        this.policy = policy;
        this.isLoading = false;
        this.cdr.markForCheck();
      });
  }

  duplicatePolicy() {
    this.policiesService.duplicateAgentPolicy(this.policyId || this.policy.id)
      .subscribe(resp => {
        if (resp?.id) {
          this.notifications.success('Agent Policy Duplicated',
            `New Agent Policy Name: ${resp?.name}`);
        }
      });
  }

  ngOnDestroy() {
    this.policySubscription?.unsubscribe();
  }
}
