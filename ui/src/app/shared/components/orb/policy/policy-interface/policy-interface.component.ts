import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges } from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';

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

  constructor() {
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
  }

  ngOnInit(): void {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode) {
      this.toggleEdit(changes.editMode.currentValue, false);
      this.updateForm();
    }
  }

  updateForm() {

  }

  toggleEdit(value, notify = true) {
    this.editMode = value;
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

}
