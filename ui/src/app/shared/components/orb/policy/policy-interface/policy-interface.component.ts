import {
  AfterViewInit,
  Component,
  EventEmitter,
  Input,
  OnChanges,
  OnInit,
  Output,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import IStandaloneEditorConstructionOptions = monaco.editor.IStandaloneEditorConstructionOptions;
import { OrbService } from 'app/common/services/orb.service';
import { EditorComponent } from 'ngx-monaco-editor';


@Component({
  selector: 'ngx-policy-interface',
  templateUrl: './policy-interface.component.html',
  styleUrls: ['./policy-interface.component.scss'],
})
export class PolicyInterfaceComponent implements OnInit, AfterViewInit, OnChanges {
  @Input()
  policy: AgentPolicy = {};

  @Input()
  editMode: boolean;

  @Output()
  editModeChange: EventEmitter<boolean>;

  @Input()
  detailsEditMode: boolean;

  @ViewChild(EditorComponent, { static: true })
  editorComponent: EditorComponent;

  editorOptions: IStandaloneEditorConstructionOptions = {
    theme: 'vs-dark',
    dragAndDrop: true,
    wordWrap: 'on',
    detectIndentation: true,
    tabSize: 2,
    autoIndent: 'full',
    trimAutoWhitespace: true,
    formatOnType: true,
    matchBrackets: 'always',
    language: 'yaml',
    automaticLayout: true,
    glyphMargin: false,
    folding: true,
    readOnly: true,
    scrollBeyondLastLine: false,
    // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
  };

  code;

  formControl: FormControl;

  constructor(
    private fb: FormBuilder,
    private orb: OrbService,
  ) {
    this.policy = {};
    this.code = '';
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
    this.detailsEditMode = false;
  }

  getCodeLineCount() {
    const editorInstance = this.editorComponent['_editor'];
    if (editorInstance) {
      const model = editorInstance.getModel();
      editorInstance.layout();
      return model ? model.getLineCount() : 0;
      
    }
    return 0;
  }
  
  ngOnInit(): void {
    this.code = this.policy.policy_data || JSON.stringify(this.policy.policy, null, 2);
  }

  ngAfterViewInit() {
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode && !changes?.editMode.firstChange) {
      this.toggleEdit(changes.editMode.currentValue, false);
    }
  }

  updateForm() {
    const { policy_data, policy } = this.policy;
    if (this.editMode) {
      this.code = policy_data || JSON.stringify(policy, null, 2);
      this.formControl = this.fb.control(this.code, [Validators.required]);
    } else {
      this.formControl = this.fb.control(null, [Validators.required]);
      this.code = policy_data || JSON.stringify(policy, null, 2);
    }
  }

  toggleEdit(edit, notify = true) {
    this.editMode = edit;
    if (this.editMode || this.detailsEditMode) {
      this.orb.pausePolling();
    }
    else {
      this.orb.startPolling();
    }
    this.editorOptions = { ...this.editorOptions, readOnly: !edit };
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    const reader: FileReader = new FileReader();
  
    reader.onload = (e: any) => {
    const fileContent = e.target.result;
    this.code = fileContent;
    };
  
    reader.readAsText(file);
  }
}
