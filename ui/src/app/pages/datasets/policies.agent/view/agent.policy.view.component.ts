import {
  ChangeDetectorRef,
  Component,
  OnChanges,
  OnDestroy,
  OnInit,
  ViewChild,
} from '@angular/core';
import {
  ActivatedRoute,
  NavigationEnd,
  Router,
  RouterEvent,
} from '@angular/router';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { PolicyConfig } from 'app/common/interfaces/orb/policy/config/policy.config.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { CodeEditorService } from 'app/common/services/code.editor.service';
import { PolicyDetailsComponent } from 'app/shared/components/orb/policy/policy-details/policy-details.component';
import { PolicyInterfaceComponent } from 'app/shared/components/orb/policy/policy-interface/policy-interface.component';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';
import yaml from 'js-yaml';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { PolicyDuplicateComponent } from '../duplicate/agent.policy.duplicate.confirmation';
import { NbDialogService } from '@nebular/theme';
import { updateMenuItems } from 'app/pages/pages-menu';
import { AgentPolicyDeleteComponent } from '../delete/agent.policy.delete.component';
import { error } from 'console';

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

  datasets: Dataset[];
  groups: AgentGroup[];

  policySubscription: Subscription;

  editMode = {
    details: false,
    interface: false,
  };

  isRequesting: boolean;

  lastUpdate: Date | null = null;

  errorConfigMessage: string;

  @ViewChild(PolicyDetailsComponent) detailsComponent: PolicyDetailsComponent;

  @ViewChild(PolicyInterfaceComponent)
  interfaceComponent: PolicyInterfaceComponent;

  constructor(
    private route: ActivatedRoute,
    private policiesService: AgentPoliciesService,
    private orb: OrbService,
    private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
    private router: Router,
    private dialogService: NbDialogService,
    private editor: CodeEditorService,
  ) {
    this.isRequesting = false;
    this.errorConfigMessage = '';
  }

  ngOnInit() {
    this.fetchData();
    updateMenuItems('Policy Management');
  }

  fetchData(newPolicyId?: any) {
    this.isLoading = true;
    if (newPolicyId) {
      this.policyId = newPolicyId;
    } else {
      this.policyId = this.route.snapshot.paramMap.get('id');
    }
    this.retrievePolicy();
    this.lastUpdate = new Date();
  }


  isEditMode() {
    const resp = Object.values(this.editMode).reduce(
      (prev, cur) => prev || cur,
      false,
    );
    if (!resp) {
      this.errorConfigMessage = '';
    }
    return resp;
  }

  canSave() {
    const detailsValid = this.editMode.details
      ? this.detailsComponent?.formGroup?.status === 'VALID'
      : true;

    const config = this.interfaceComponent?.code;
    let interfaceValid = false;

    if (this.policy.format === 'json') {
      if (this.editor.isJson(config)) {
        interfaceValid = true;
        this.errorConfigMessage = '';
      }
      else {
        interfaceValid = false;
        this.errorConfigMessage = 'Invalid JSON configuration, check syntax errors';
      }
    }
    else if (this.policy.format === 'yaml') {
      if (this.editor.isYaml(config)) {
        interfaceValid = true;
        this.errorConfigMessage = '';
      }
      else {
        interfaceValid = false;
        this.errorConfigMessage = 'Invalid YAML configuration, check syntax errors';
      }
    }
    return detailsValid && interfaceValid;
  }

  discard() {
    this.editMode.details = false;
    this.editMode.interface = false;
  }

  save() {
    this.isRequesting = true;

    const { format, version, name, description, id, backend } = this.policy;

    // get values from all modified sections' forms and submit through service.
    const policyDetails = this.detailsComponent.formGroup?.value;
    const tags = this.detailsComponent.selectedTags;
    const policyInterface = this.interfaceComponent.code;

    // trying to work around rest api
    const detailsPartial = (!!this.editMode.details && {
      ...policyDetails,
    }) || { name, description };

    let interfacePartial = {};

    try {
      if (format === 'yaml') {
        if (this.editor.isJson(policyInterface)) {
          throw new Error('Invalid YAML format');
        }
        yaml.load(policyInterface);

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
        version,
        id,
        tags,
        backend,
      } as AgentPolicy;

      this.policiesService.editAgentPolicy(payload).subscribe(
        (resp) => {
        this.notifications.success('Agent Policy updated successfully', '');
        this.discard();
        this.policy = resp;
        this.orb.refreshNow();
        this.isRequesting = false;
        },
        (err) => {
          this.isRequesting = false;
        },
        );

    } catch (err) {
      this.notifications.error(
        'Failed to edit Agent Policy',
        `Error: Invalid ${format.toUpperCase()}`,
      );
      this.isRequesting = false;
    }
  }

  retrievePolicy() {
    this.policySubscription = this.orb
      .getPolicyFullView(this.policyId)
      .subscribe(({ policy, datasets, groups }) => {
        this.policy = policy;
        this.datasets = datasets;
        this.groups = groups;
        this.isLoading = false;
        this.cdr.markForCheck();
      });
  }
  onOpenDuplicatePolicy() {
    const policy = this.policy.name;
    this.dialogService
      .open(PolicyDuplicateComponent, {
        context: { policy },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.duplicatePolicy(this.policy);
        }
      });
  }
  duplicatePolicy(agentPolicy: any) {
    this.policiesService
    .duplicateAgentPolicy(agentPolicy.id)
    .subscribe((newAgentPolicy) => {
      if (newAgentPolicy?.id) {
        this.notifications.success(
          'Agent Policy Duplicated',
          `New Agent Policy Name: ${newAgentPolicy?.name}`,
        );
        this.router.navigateByUrl(`/pages/datasets/policies/view/${newAgentPolicy?.id}`);
        this.fetchData(newAgentPolicy.id);
      }
    });
  }

  ngOnDestroy() {
    this.policySubscription?.unsubscribe();
    this.orb.isPollingPaused ? this.orb.startPolling() : null;
    this.orb.killPolling.next();
  }
  openDeleteModal() {
    const { name: name, id } = this.policy as AgentPolicy;
    this.dialogService
      .open(AgentPolicyDeleteComponent, {
        context: { name },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.policiesService.deleteAgentPolicy(id).subscribe(() => {
            this.notifications.success(
              'Agent Policy successfully deleted',
              '',
            );
            this.goBack();
          });
        }
      });
  }
  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  hasChanges() {
    const policyDetails = this.detailsComponent.formGroup?.value;
    const tags = this.detailsComponent.selectedTags;

    const description = this.policy.description ? this.policy.description : '';
    const formsDescription = policyDetails.description === null ? '' : policyDetails.description;

    const selectedTags = JSON.stringify(tags);
    const orb_tags = JSON.stringify(this.policy.tags);

    if (policyDetails.name !== this.policy.name || formsDescription !== description || selectedTags !== orb_tags) {
      return true;
    }
    return false;
  }
}
