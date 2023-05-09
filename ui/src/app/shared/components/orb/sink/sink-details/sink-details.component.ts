import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { Tags } from 'app/common/interfaces/orb/tag';


@Component({
  selector: 'ngx-sink-details',
  templateUrl: './sink-details.component.html',
  styleUrls: ['./sink-details.component.scss']
})
export class SinkDetailsComponent implements OnInit, OnChanges {
  @Input()
  sink: Sink;

  @Input()
  editMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  formGroup: FormGroup;

  selectedTags: Tags;

  constructor(private fb: FormBuilder) { 
    this.sink = {};
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
  }

  ngOnInit(): void {
    this.selectedTags = this.sink?.tags || {};
  }

  ngOnChanges(changes: SimpleChanges): void {
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
          ],
        ],
        description: [description],
      });
      this.selectedTags = {...tags} || {};
    } else {
      this.formGroup = this.fb.group({
        name: null,
        description: null,
      });
    }
  }

  toggleEdit(value, notify = true) {
    this.editMode = value;
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }
}
