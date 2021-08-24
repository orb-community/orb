import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

const strings = STRINGS.sink;

@Component({
  selector: 'ngx-sinks-details-component',
  templateUrl: './sinks.details.component.html',
  styleUrls: ['./sinks.details.component.scss'],
})
export class SinksDetailsComponent {
  header = strings.details.header;
  close = strings.details.close;
  name = strings.propNames.name;
  description = strings.propNames.description;
  backend = strings.propNames.backend;
  remote_host = strings.propNames.config_remote_host;
  ts_created = strings.propNames.ts_created;

  @Input() sink: Sink = {};

  constructor(
    protected dialogRef: NbDialogRef<SinksDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }


  onOpenEdit(row: any) {
    this.router.navigate(['../sinks/edit'], {
      relativeTo: this.route,
      queryParams: {id: row.id},
      state: {sink: row},
    });
  }

  onClose() {
    this.dialogRef.close();
  }
}
