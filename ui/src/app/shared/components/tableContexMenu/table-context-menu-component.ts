import { Component, Input } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { NbDialogService } from '@nebular/theme';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { OrbService } from 'app/common/services/orb.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { AgentPolicyDeleteComponent } from 'app/pages/datasets/policies.agent/delete/agent.policy.delete.component';
import { PolicyDuplicateComponent } from 'app/pages/datasets/policies.agent/duplicate/agent.policy.duplicate.confirmation';
import { AgentDeleteComponent } from 'app/pages/fleet/agents/delete/agent.delete.component';
import { AgentGroupDeleteComponent } from 'app/pages/fleet/groups/delete/agent.group.delete.component';
import { AgentGroupDetailsComponent } from 'app/pages/fleet/groups/details/agent.group.details.component';
import { SinkDeleteComponent } from 'app/pages/sinks/delete/sink.delete.component';

@Component({
  selector: 'ngx-table-context-menu',
  templateUrl: './table-context-menu-component.html',
  styleUrls: ['./table-context-menu-component.scss']
})
export class TableContextMenu {

  @Input()
  items: any[];

  @Input()
  top: number;

  @Input()
  left: number;

  @Input()
  rowObject: any;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private dialogService: NbDialogService,
    private orb: OrbService,
    protected agentService: AgentsService,
    protected notificationsService: NotificationsService,
    private agentGroupsService: AgentGroupsService,
    private sinkService: SinksService,
    private agentPoliciesService: AgentPoliciesService,
  ) {

  }
  handleClick(item: any) {
    const { action } = item;
    if (action === 'openview') {
      this.openView();
    } else if (action === 'opendelete') {
      this.openDelete();
    } else if (action === 'openduplicate') {
      this.openDuplicate();
    } else if (action === 'openedit') {
      this.openGroupEdit();
    }
  }

  openGroupEdit() {
    this.router.navigate([`edit/${this.rowObject.id}`], {
      state: { agentGroup: this.rowObject, edit: true },
      relativeTo: this.route,
    });
  }

  openView() {
    const { objectType, id } = this.rowObject;
    if (objectType === 'group') {
      this.dialogService
      .open(AgentGroupDetailsComponent, {
        context: { agentGroup: this.rowObject },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((resp) => {
        if (resp) {
          this.openGroupEdit();
        }
      });
    }
    else {
      this.router.navigate([`view/${id}`], {
        relativeTo: this.route,
      });
    }
  }
  openDelete() {
    const { objectType, name, id } = this.rowObject;
  
    const deleteCallback = () => {
      this.notificationsService.success(
        `${objectType.charAt(0).toUpperCase() + objectType.slice(1)} successfully deleted`,
        ''
      );
      this.orb.refreshNow();
    };
  
    if (objectType === 'agent') {
      this.dialogService
        .open(AgentDeleteComponent, {
          context: { name },
          autoFocus: true,
          closeOnEsc: true,
        })
        .onClose.subscribe((confirm) => {
          if (confirm) {
            this.agentService.deleteAgent(id).subscribe(deleteCallback);
          }
        });
    } else if (objectType === 'group') {
      this.dialogService
        .open(AgentGroupDeleteComponent, {
          context: { name },
          autoFocus: true,
          closeOnEsc: true,
        })
        .onClose.subscribe((confirm) => {
          if (confirm) {
            this.agentGroupsService.deleteAgentGroup(id).subscribe(deleteCallback);
          }
        });
    } else if (objectType === 'policy') {
      this.dialogService
        .open(AgentPolicyDeleteComponent, {
          context: { name },
          autoFocus: true,
          closeOnEsc: true,
        })
        .onClose.subscribe((confirm) => {
          if (confirm) {
            this.agentPoliciesService.deleteAgentPolicy(id).subscribe(deleteCallback);
          }
        });
    } else if (objectType === 'sink') {
      this.dialogService
        .open(SinkDeleteComponent, {
          context: { sink: this.rowObject },
          autoFocus: true,
          closeOnEsc: true,
        })
        .onClose.subscribe((confirm) => {
          if (confirm) {
            this.sinkService.deleteSink(id).subscribe(deleteCallback);
          }
        });
    }
  }

  openDuplicate() {
    const policy = this.rowObject.name;
    this.dialogService
      .open(PolicyDuplicateComponent, {
        context: { policy },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.duplicatePolicy(this.rowObject);
        }
      });
  }
  duplicatePolicy(agentPolicy: any) {
    this.agentPoliciesService
    .duplicateAgentPolicy(agentPolicy.id)
    .subscribe((newAgentPolicy) => {
      if (newAgentPolicy?.id) {
        this.notificationsService.success(
          'Agent Policy Duplicated',
          `New Agent Policy Name: ${newAgentPolicy?.name}`,
        );
        this.router.navigate([`view/${newAgentPolicy.id}`], {
          relativeTo: this.route,
        });
      }
    });
  }
}
