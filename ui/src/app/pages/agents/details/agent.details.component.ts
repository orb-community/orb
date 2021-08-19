import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';

const strings = STRINGS.sink;

@Component({
  selector: 'ngx-agent-details-component',
  templateUrl: './agent.details.component.html',
  styleUrls: ['./agent.details.component.scss'],
})
export class AgentDetailsComponent {
  header = strings.details.header;
  close = strings.details.close;
  name = strings.propNames.name;
  description = strings.propNames.description;
  backend = strings.propNames.backend;
  remote_host = strings.propNames.config_remote_host;
  ts_created = strings.propNames.ts_created;

  @Input() agent = {
    id: '',
    name: '',
    description: '',
    backend: '',
    config: {
      remote_host: '',
      username: '',
    },
    ts_created: '',
  };

  constructor(
    protected dialogRef: NbDialogRef<AgentDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }


  onOpenEdit(row: any) {
    this.router.navigate(['../agents/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {sink: row},
    });
  }

  onClose() {
    this.dialogRef.close();
  }
}
