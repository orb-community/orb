import { Component, Input } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'ngx-sinks-add-component',
  templateUrl: './sinks.add.component.html',
  styleUrls: ['./sinks.add.component.scss'],
})
export class SinksAddComponent {
  editorMetadata = '';

  @Input() formData = {
    name: '',
    description: '',
    tags: '',
    backend: '',
    config: {
      remote_host: '',
      username: '',
    },
    metadata: {},
  };

  constructor(
    private sinksService: SinksService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
  }

  goBack() {
    this.router.navigate(['../../sinks'], {relativeTo: this.route});
  }

  cancel() {
    this.goBack();
  }

  submit() {
    const action = this.router.routerState.snapshot.url.split('/').pop().toLowerCase();

    if (this.formData.tags !== '') {
      try {
        this.formData.tags = JSON.parse(this.formData.tags);
      } catch (e) {
        this.notificationsService.error('Wrong metadata format', '');
        return;
      }
    }

    this.formData.backend && (this.formData.metadata['backend'] = this.formData.backend);
    if (action === 'add') {
      this.sinksService.addSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink successfully created', '');
        },
        error => {
          this.notificationsService.error('Sink creation failed', error);
        },
        () => {
          this.goBack();
        },
      );
    }
    if (action === 'edit') {
      this.sinksService.editSink(this.formData).subscribe(
        resp => {
          this.notificationsService.success('Sink successfully updated', '');
        },
        error => {
          this.notificationsService.error('Sink update failed', error);
        },
        () => {
          this.goBack();
        },
      );
    }
  }
}
