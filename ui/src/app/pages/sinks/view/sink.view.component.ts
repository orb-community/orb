import { ChangeDetectorRef, Component, OnChanges, OnDestroy, OnInit, SimpleChanges, ViewChild } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkConfigComponent } from 'app/shared/components/orb/sink/sink-config/sink-config.component';
import { SinkDetailsComponent } from 'app/shared/components/orb/sink/sink-details/sink-details.component';
import { STRINGS } from 'assets/text/strings';
import { Subscription } from 'rxjs';
import { updateMenuItems } from 'app/pages/pages-menu';

@Component({
  selector: 'ngx-sink-view',
  templateUrl: './sink.view.component.html',
  styleUrls: ['./sink.view.component.scss']
})
export class SinkViewComponent implements OnInit, OnChanges, OnDestroy {
  strings = STRINGS;
  
  isLoading = false;

  sink: Sink;

  sinkId = '';

  sinkSubscription: Subscription;

  editMode = {
    details: false,
    config: false,
  }

  @ViewChild(SinkDetailsComponent) detailsComponent: SinkDetailsComponent;

  @ViewChild(SinkConfigComponent)
  configComponent: SinkConfigComponent;

  constructor(private cdr: ChangeDetectorRef,
    private notifications: NotificationsService,
    private sinks: SinksService,
    private route: ActivatedRoute,
    ) { }

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
    const detailsValid = this.editMode.details
      ? this.detailsComponent?.formGroup?.status === 'VALID'
      : true;

    const configValid = this.editMode.config
      ? this.configComponent?.formControl?.status === 'VALID'
      : true;

    return detailsValid && configValid;
  }

  discard() {
    this.editMode.details = false;
    this.editMode.config = false;
  }

save() {
  const { id, backend } = this.sink;
  const sinkDetails = this.detailsComponent.formGroup?.value;
  const tags = this.detailsComponent.selectedTags;
  const configSink = this.configComponent.code;

  const details = { ...sinkDetails, tags };
  const isJson = this.isJson(configSink);

  let payload: Sink = { id, backend };

  if (isJson) {
    const config = JSON.parse(configSink);

    if (this.editMode.details && !this.editMode.config) {
      payload = { ...payload, ...details };
    } else if (!this.editMode.details && this.editMode.config) {
      payload = { ...payload, config };
    } else {
      payload = { ...payload, ...details, config };
    }
  } else {
    if (this.editMode.details && !this.editMode.config) {
      payload = { ...payload, ...details };
    } else if (!this.editMode.details && this.editMode.config) {
      payload = { ...payload, format: 'yaml', config_data: configSink };
    } else {
      payload = { ...payload, ...details, format: 'yaml', config_data: configSink };
    }
  }

  try {
    this.sinks.editSink(payload).subscribe((resp) => {
      this.discard();
      this.sink = resp;
      this.fetchData();
      this.notifications.success('Sink updated successfully', '');
    });
  } catch (err) {
    this.notifications.error('Failed to edit Sink', 'Error: Invalid configuration');
  }
}
  isJson(str: string) {
    try {
        JSON.parse(str);
        return true;
    } catch {
        return false;
    }
}
  retrieveSink() {
    this.sinkSubscription = this.sinks
    .getSinkById(this.sinkId)
    .subscribe(sink => {
      this.sink = sink;
      this.isLoading = false;
      this.cdr.markForCheck();
    });
  }

  ngOnDestroy(): void {
    this.sinkSubscription.unsubscribe();
  }
}
