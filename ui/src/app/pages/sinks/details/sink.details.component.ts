import { Component, Input, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink, SinkBackends, SinkStates } from 'app/common/interfaces/orb/sink.interface';

@Component({
  selector: 'ngx-sink-details-component',
  templateUrl: './sink.details.component.html',
  styleUrls: ['./sink.details.component.scss'],
})
export class SinkDetailsComponent implements OnInit {

  exporterField: string;

  sinkBackends = SinkBackends;

  strings = STRINGS.sink;

  @Input() sink: Sink = {};

  _sinkStates = SinkStates;

  constructor(
    protected dialogRef: NbDialogRef<SinkDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    !this.sink.tags ? this.sink.tags = {} : null;
    this.exporterField = '';
  }

  onOpenEdit(sink: any) {
    this.router.navigateByUrl(`/pages/sinks/edit/${sink.id}`);
    this.dialogRef.close();
  }

  onClose() {
    this.dialogRef.close(false);
  }

  onOpenView(sink: any) {
    this.router.navigateByUrl(`/pages/sinks/view/${sink.id}`);
    this.dialogRef.close();
  }
  ngOnInit() {
    const exporter = this.sink.config.exporter;
    this.exporterField = exporter.remote_host !== undefined ? 'Remote Host URL' : 'Endpoint URL';
  }
}
