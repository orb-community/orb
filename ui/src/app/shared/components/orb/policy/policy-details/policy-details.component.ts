import {
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
} from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'ngx-policy-details',
  templateUrl: './policy-details.component.html',
  styleUrls: ['./policy-details.component.scss'],
})
export class PolicyDetailsComponent implements OnInit, OnChanges {
  @Input()
  policy: AgentPolicy;

  @Input()
  editMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  formGroup: FormGroup;

  constructor(private fb: FormBuilder) {
    this.policy = {};
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
  }

  ngOnInit(): void {}

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode) {
      this.toggleEdit(changes.editMode.currentValue, false);
    }
  }

  updateForm() {
    if (this.editMode) {
      const { name: name, description } = this.policy;
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
