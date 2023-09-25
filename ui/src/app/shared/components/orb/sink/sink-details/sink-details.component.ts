import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Sink, SinkBackends, SinkStates } from 'app/common/interfaces/orb/sink.interface';
import { SinkFeature } from 'app/common/interfaces/orb/sink/sink.feature.interface';
import { Tags } from 'app/common/interfaces/orb/tag';
import { SinksService } from 'app/common/services/sinks/sinks.service';


@Component({
  selector: 'ngx-sink-details',
  templateUrl: './sink-details.component.html',
  styleUrls: ['./sink-details.component.scss']
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

  formGroup: FormGroup;

  selectedTags: Tags;

  mode: string;

  sinkTypesList = [];

  sinkStates = SinkStates;

  constructor(
    private fb: FormBuilder,
    private sinksService: SinksService,
    ) { 
    this.sink = {};
    this.createMode = false;
    this.editMode = false;
    this.mode = 'read';
    this.sinkBackend = new EventEmitter<string>();
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
    Promise.all([this.getSinkBackends()]).then((responses) => {
      const backends = responses[0];
      this.sinkTypesList = backends.map(entry => entry.backend);
    })
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
    }
    else if (this.createMode) {

      const { name, description, backend, tags } = this.sink;
      
      this.formGroup = this.fb.group({
        name: [name, [Validators.required, Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$'), Validators.maxLength(64)]],
        description: [description, [Validators.maxLength(64)]],
        backend: [backend, Validators.required],
      });

      this.selectedTags = { ...tags };
    } 
    else {
      this.formGroup = this.fb.group({
        name: null,
        description: null,
        backend: null,
      });
    }
  }

  toggleEdit(value, notify = true) {
    this.editMode = value;
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

  getMode() {
    if(this.editMode == true) {
      this.mode = 'edit';
    }
    else if (this.createMode == true) {
      this.mode = 'create';
    }
    else {
      this.mode = 'read';
    }
  }
  
  getSinkBackends() {
    return new Promise<SinkFeature[]>(resolve => {
      this.sinksService.getSinkBackends().subscribe(backends => {
        resolve(backends);
      });
    });
  }

  onChangeConfigBackend(backend: any) {
    this.sinkBackend.emit(backend);
  }
}
