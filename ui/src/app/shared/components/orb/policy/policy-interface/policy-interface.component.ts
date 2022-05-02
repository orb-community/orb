import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { FormBuilder, FormControl, Validators } from '@angular/forms';

@Component({
  selector: 'ngx-policy-interface',
  templateUrl: './policy-interface.component.html',
  styleUrls: ['./policy-interface.component.scss'],
})
export class PolicyInterfaceComponent implements OnInit, OnChanges {
  @Input()
  policy: AgentPolicy = {};

  @Input()
  editMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  formControl: FormControl;

  constructor(
    private fb: FormBuilder,
  ) {
    this.policy = {};
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode) {
      this.toggleEdit(changes.editMode.currentValue, false);
    }
  }

  updateForm() {
    if (this.editMode) {
      const { policy_data } = this.policy;
      this.formControl = this.fb.control(policy_data, [Validators.required]);
    } else {
      this.formControl = this.fb.control(null, [Validators.required]);
    }
  }

  toggleEdit(value, notify = true) {
    this.editMode = value;
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

}
