import { ChangeDetectorRef, Component, OnChanges, OnDestroy, OnInit, SimpleChanges, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { Sink, SinkStates } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkConfigComponent } from 'app/shared/components/orb/sink/sink-config/sink-config.component';
import { SinkDetailsComponent } from 'app/shared/components/orb/sink/sink-details/sink-details.component';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';
import { updateMenuItems } from 'app/pages/pages-menu';
import * as YAML from 'yaml';
import { CodeEditorService } from 'app/common/services/code.editor.service';
import { SinkDeleteComponent } from '../delete/sink.delete.component';
import { NbDialogService } from '@nebular/theme';
import { OrbService } from 'app/common/services/orb.service';

@Component({
  selector: 'ngx-sink-view',
  templateUrl: './sink.view.component.html',
  styleUrls: ['./sink.view.component.scss'],
})
export class SinkViewComponent implements OnInit, OnChanges, OnDestroy {
  strings = STRINGS;

  isLoading = false;

  sink: Sink;

  sinkId = '';

  sinkSubscription: Subscription;

  lastUpdate: Date | null = null;

  sinkStates = SinkStates;

  editMode = {
    details: false,
    config: false,
  };

  isRequesting: boolean;

  @ViewChild(SinkDetailsComponent) detailsComponent: SinkDetailsComponent;

  @ViewChild(SinkConfigComponent)
  configComponent: SinkConfigComponent;

  constructor(private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
    private sinks: SinksService,
    private route: ActivatedRoute,
    private editor: CodeEditorService,
    private dialogService: NbDialogService,
    private router: Router,
    private orb: OrbService,
    ) {
      this.isRequesting = false;
    }

  ngOnInit(): void {
    this.fetchData();
    updateMenuItems('Sink Management');
  }

  ngOnChanges(): void {
    this.fetchData();
  }

  fetchData() {
    this.isLoading = true;
    this.sinkId = this.route.snapshot.paramMap.get('id');
    this.retrieveSink();
  }

  isEditMode() {
    return Object.values(this.editMode).reduce(
      (prev, cur) => prev || cur,
      false,
    );
  }

  canSave() {
    let configValid = true;
    const detailsValid = this.editMode.details
      ? this.detailsComponent?.formGroup?.status === 'VALID'
      : true;

    const configSink = this.configComponent?.code;
    let config;

    if (this.editor.isJson(configSink)) {
      config = JSON.parse(configSink);
    } else if (this.editor.isYaml(configSink)) {
      config = YAML.parse(configSink);
    } else {
        return false;
    }

    if (this.editMode.config) {
      configValid = !this.editor.checkEmpty(config.authentication) && !this.editor.checkEmpty(config.exporter) && !this.checkString(config);
    }
    return detailsValid && configValid;
  }

  checkString(config: any): boolean {
    if (typeof config.authentication.password !== 'string' || typeof config.authentication.username !== 'string') {
      return true;
    }
    return false;
  }

  discard() {
    this.editMode.details = false;
    this.editMode.config = false;
  }

  save() {
    this.isRequesting = true;
    const { id, backend } = this.sink;
    const sinkDetails = this.detailsComponent.formGroup?.value;
    const tags = this.detailsComponent.selectedTags;
    const configSink = this.configComponent.code;

    const details = { ...sinkDetails, tags };

    try {
      let payload: any;
      if (this.editMode.config && !this.editMode.details) {
        payload = { id, backend, config: {}};
        const config = YAML.parse(configSink);
        payload.config = config;
      }
      if (this.editMode.details && !this.editMode.config) {
        payload = { id, backend, ...details };
      }
      if (this.editMode.details && this.editMode.config) {
        payload = { id, backend, ...details, config: {}};
        const config = YAML.parse(configSink);
        payload.config = config;
      }
      this.sinks.editSink(payload as Sink).subscribe(
        (resp) => {
          this.discard();
          this.sink = resp;
          this.orb.refreshNow();
          this.notifications.success('Sink updated successfully', '');
          this.isRequesting = false;
        },
        (err) => {
          this.isRequesting = false;
        });
    } catch (err) {
      this.notifications.error('Failed to edit Sink', 'Error: Invalid configuration');
    }
  }

  retrieveSink() {
    this.sinkSubscription = this.orb.getSinkView(this.sinkId).subscribe(sink => {
      this.sink = sink;
      this.isLoading = false;
      this.cdr.markForCheck();
      this.lastUpdate = new Date();
    });
  }

  ngOnDestroy(): void {
    this.sinkSubscription.unsubscribe();
    this.orb.isPollingPaused ? this.orb.startPolling() : null;
    this.orb.killPolling.next();
  }
  openDeleteModal() {
    const { id } = this.sink;
    this.dialogService
      .open(SinkDeleteComponent, {
        context: { sink: this.sink },
        autoFocus: true,
        closeOnEsc: true,
      })
      .onClose.subscribe((confirm) => {
        if (confirm) {
          this.sinks.deleteSink(id).subscribe(() => {
            this.notifications.success('Sink successfully deleted', '');
            this.goBack();
          });
        }
      });
  }
  goBack() {
    this.router.navigateByUrl('/pages/sinks');
  }
  hasChanges() {
    const sinkDetails = this.detailsComponent.formGroup?.value;
    const tags = this.detailsComponent.selectedTags;
    const selectedTags = JSON.stringify(tags);
    const orb_tags = this.sink.tags ? JSON.stringify(this.sink.tags) : '{}';

    if (sinkDetails.name !== this.sink.name || sinkDetails?.description !== this.sink?.description || selectedTags !== orb_tags) {
      return true;
    }
    return false;
  }
}
