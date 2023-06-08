import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { Sink, SinkBackends } from 'app/common/interfaces/orb/sink.interface';
import IStandaloneEditorConstructionOptions = monaco.editor.IStandaloneEditorConstructionOptions;
@Component({
  selector: 'ngx-sink-config',
  templateUrl: './sink-config.component.html',
  styleUrls: ['./sink-config.component.scss']
})
export class SinkConfigComponent implements OnInit, OnChanges {

  @Input()
  sink: Sink;

  @Input()
  editMode: boolean;

  @Input()
  createMode: boolean;

  @Input()
  sinkBackend: string;
  
  @Output()
  editModeChange: EventEmitter<boolean>;

  @ViewChild('editorComponent')
  editor;

  editorOptions: IStandaloneEditorConstructionOptions = {
    theme: 'vs-dark',
    dragAndDrop: true,
    wordWrap: 'on',
    detectIndentation: true,
    tabSize: 2,
    autoIndent: 'full',
    formatOnPaste: true,
    trimAutoWhitespace: true,
    formatOnType: true,
    matchBrackets: 'always',
    language: 'json',
    automaticLayout: true,
    glyphMargin: false,
    folding: true,
    readOnly: true,
    scrollBeyondLastLine: false,
    // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
  };

  code = '';

  sinkConfigSchemaPrometheus: any;

  sinkConfigSchemaOtlp: any;

  formControl: FormControl;

  constructor(private fb: FormBuilder) { 
    this.sink = {};
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.updateForm();
    this.sinkConfigSchemaPrometheus = {
      "authentication" : {
        "type": "basicauth", 
        "password": "",
        "username": "",
      },
      "exporter" : {
        "remote_host": "",
      },
      "opentelemetry": "enabled",
    }
    this.sinkConfigSchemaOtlp = {
      "authentication" : {
        "type": "basicauth", 
        "password": "",
        "username": "",
      },
      "exporter" : {
        "endpoint": "",
      },
      "opentelemetry": "enabled",
    }
  }

  ngOnInit(): void {
    if (this.createMode) {
      this.toggleEdit(true);
      this.code = JSON.stringify(this.sinkConfigSchemaOtlp, null, 2);
    }
    else {
      this.code = JSON.stringify(this.sink.config, null, 2);
    }
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes?.editMode && !changes?.editMode.firstChange) {
      this.toggleEdit(changes.editMode.currentValue, false);
    }
    if (changes?.sinkBackend) {
      if (this.sinkBackend === SinkBackends.prometheus) {
        this.code = JSON.stringify(this.sinkConfigSchemaPrometheus, null, 2);
      }
      else {
        this.code = JSON.stringify(this.sinkConfigSchemaOtlp, null, 2);
      }
    }
  }

  updateForm() {
    const { config } = this.sink;
    if (this.editMode) {
      this.code = JSON.stringify(config, null, 2);
      this.formControl = this.fb.control(this.code, [Validators.required]);
    } else {
      this.formControl = this.fb.control(null, [Validators.required]);
      this.code = JSON.stringify(config, null, 2);
    }
  }

  toggleEdit(edit, notify = true) {
    this.editMode = edit;
    this.editorOptions = { ...this.editorOptions, readOnly: !edit };
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }

}
