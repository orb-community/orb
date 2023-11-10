import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Sink, SinkBackends, SinkStates } from 'app/common/interfaces/orb/sink.interface';
import { SinkFeature } from 'app/common/interfaces/orb/sink/sink.feature.interface';
import { Tags } from 'app/common/interfaces/orb/tag';
import { OrbService } from 'app/common/services/orb.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';


@Component({
  selector: 'ngx-sink-details',
  templateUrl: './sink-details.component.html',
  styleUrls: ['./sink-details.component.scss'],
})

export class SinkDetailsComponent implements OnInit, OnChanges {

  @Output()
  sinkBackend: EventEmitter<string>;

  @Input()
  sink: Sink;

  @Input()
  editMode: boolean;

  @Input()
  createMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  @Input()
  configEditMode: boolean;

  @Input()
  sinkTypesList = [];

  formGroup: FormGroup;

  selectedTags: Tags;

  mode: string;

  sinkStates = SinkStates;

  constructor(
    private fb: FormBuilder,
    private sinksService: SinksService,
    private orb: OrbService,
    ) {
    this.sink = {};
    this.createMode = false;
    this.editMode = false;
    this.mode = 'read';
    this.sinkBackend = new EventEmitter<string>();
    this.editModeChange = new EventEmitter<boolean>();
    this.configEditMode = false;
    this.updateForm();
  }

  ngOnInit(): void {
    this.getMode();
    this.selectedTags = this.sink?.tags || {};
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.getMode();
    if (changes?.editMode) {
      this.toggleEdit(changes.editMode.currentValue, false);
    }
    if (changes?.sink) {
      this.selectedTags = this.sink?.tags || {};
    }
  }

  updateForm() {
    if (this.editMode) {
      const { name, description, tags } = this.sink;
      this.formGroup = this.fb.group({
        name: [
          name,
          [
            Validators.required,
            Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$'),
            Validators.maxLength(64),
          ],
        ],
        description: [description],
      });
      this.selectedTags = {...tags} || {};
    } else if (this.createMode) {

      const { name, description, backend, tags } = this.sink;

      this.formGroup = this.fb.group({
        name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$'), Validators.maxLength(64)]],
        description: [description, [Validators.maxLength(64)]],
        backend: [backend, Validators.required],
      });

      this.selectedTags = { ...tags };
    } else {
      this.formGroup = this.fb.group({
        name: null,
        description: null,
        backend: null,
      });
    }
  }

  toggleEdit(value, notify = true) {
    this.editMode = value;
    if (this.editMode || this.configEditMode) {
      this.orb.pausePolling();
    } else {
      this.orb.startPolling();
    }
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

  getMode() {
    if (this.editMode === true) {
      this.mode = 'edit';
    } else if (this.createMode === true) {
      this.mode = 'create';
    } else {
      this.mode = 'read';
    }
  }

  onChangeConfigBackend(backend: any) {
    this.sinkBackend.emit(backend);
  }
}
