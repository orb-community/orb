import { Component, Input } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { STRINGS } from 'assets/text/strings';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';

@Component({
  selector: 'ngx-dataset-details-component',
  templateUrl: './dataset.details.component.html',
  styleUrls: ['./dataset.details.component.scss'],
})
export class DatasetDetailsComponent {
  @Input() dataset: Dataset = {};

  constructor(
    protected dialogRef: NbDialogRef<DatasetDetailsComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
  ) {
    !this.dataset.tags ? this.dataset.tags = {} : null;
  }

  onOpenEdit(dataset: any) {
    this.dialogRef.close(true);
  }

  onClose() {
    this.dialogRef.close(false);
  }
}
