import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';

const strings = STRINGS.sink;

@Component({
  selector: 'ngx-sink-details-component',
  templateUrl: './sink.details.component.html',
  styleUrls: ['./sink.details.component.scss'],
})
export class SinkDetailsComponent {
  header = strings.details.header;
  close = strings.details.close;
  name = strings.propNames.name;
  description = strings.propNames.description;
  backend = strings.propNames.backend;
  remote_host = strings.propNames.config_remote_host;
  ts_created = strings.propNames.ts_created;

  @Input() sink: Sink = {};

  constructor(
    protected dialogRef: NbDialogRef<SinkDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
  }

  onOpenEdit(sink: any) {
    this.router.navigate(
      [`../sink/edit/${sink.id}`, sink.id], {
        relativeTo: this.route,
      });
  }

  onClose() {
    this.dialogRef.close();
  }
}
