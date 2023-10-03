import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { Sink, SinkBackends } from 'app/common/interfaces/orb/sink.interface';
import * as YAML from 'yaml';
import IStandaloneEditorConstructionOptions = monaco.editor.IStandaloneEditorConstructionOptions;
import { OrbService } from 'app/common/services/orb.service';

@Component({
  selector: 'ngx-sink-config',
  templateUrl: './sink-config.component.html',
  styleUrls: ['./sink-config.component.scss'],
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

  @Input()
  detailsEditMode: boolean;

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
  editorOptionsYaml = {
    theme: 'vs-dark',
    language: 'yaml',
    automaticLayout: true,
    glyphMargin: false,
    folding: true,
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
    readOnly: true,
  };
  code = '';

  sinkConfigSchemaPrometheus: any;

  sinkConfigSchemaOtlp: any;

  formControl: FormControl;

  isYaml: boolean;

  constructor(
    private fb: FormBuilder,
    private orb: OrbService,
    ) {
    this.isYaml = true;
    this.sink = {};
    this.editMode = false;
    this.editModeChange = new EventEmitter<boolean>();
    this.detailsEditMode = false;
    this.updateForm();
    this.sinkConfigSchemaPrometheus = {
      'authentication' : {
        'type': 'basicauth',
        'password': '',
        'username': '',
      },
      'exporter' : {
        'remote_host': '',
      },
      'opentelemetry': 'enabled',
    };
    this.sinkConfigSchemaOtlp = {
      'authentication' : {
        'type': 'basicauth',
        'password': '',
        'username': '',
      },
      'exporter' : {
        'endpoint': '',
      },
      'opentelemetry': 'enabled',
    };
  }

  ngOnInit(): void {
    if (this.createMode) {
      this.toggleEdit(true);
      this.code = YAML.stringify(this.sinkConfigSchemaOtlp);
    } else {
      // if (this.sink.config_data && this.sink.format === 'yaml') {
      // this.isYaml = true;
      const parsedCode = YAML.parse(JSON.stringify(this.sink.config));
      this.code = YAML.stringify(parsedCode);
      // }
      // else if (this.isJson(JSON.stringify(this.sink.config))) {
      //   this.isYaml = false;
      //   this.code = JSON.stringify(this.sink.config, null, 2);
      // }

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
ngOnChanges(changes: SimpleChanges) {
  const { editMode, sinkBackend } = changes;
  if (editMode && !editMode.firstChange) {
    this.toggleEdit(editMode.currentValue, false);
  }
  if (sinkBackend) {
    const sinkConfigSchema = this.sinkBackend === SinkBackends.prometheus
      ? this.sinkConfigSchemaPrometheus
      : this.sinkConfigSchemaOtlp;

    this.code = this.isYaml
      ? YAML.stringify(sinkConfigSchema, null)
      : JSON.stringify(sinkConfigSchema, null, 2);
    this.code = YAML.stringify(sinkConfigSchema, null);
  }
}

updateForm() {
  const configData = this.sink.config;
  // const isYamlFormat = this.sink.format === 'yaml';

  if (this.editMode) {
    // this.isYaml = isYamlFormat;
    // this.code = isYamlFormat ? YAML.stringify(configData) : JSON.stringify(this.sink.config, null, 2);
    this.code = YAML.stringify(configData);
  } else {
    this.formControl = this.fb.control(null, [Validators.required]);
    // this.isYaml = isYamlFormat;
    // this.code = isYamlFormat ? YAML.stringify(configData) : JSON.stringify(this.sink.config, null, 2);
    this.code = YAML.stringify(configData);
  }

  this.formControl = this.fb.control(this.code, [Validators.required]);
}

  toggleEdit(edit, notify = true) {
    this.editMode = edit;
    if ((this.editMode || this.detailsEditMode) && !this.createMode) {
      this.orb.pausePolling();
    } else {
      this.orb.startPolling();
    }
    this.editorOptions = { ...this.editorOptions, readOnly: !edit };
    this.editorOptionsYaml = { ...this.editorOptionsYaml, readOnly: !edit };
    this.updateForm();
    !!notify && this.editModeChange.emit(this.editMode);
  }
  toggleLanguage() {
    this.isYaml = !this.isYaml;
    if (this.isYaml) {
      const parsedCode = YAML.parse(this.code);
      this.code = YAML.stringify(parsedCode);
    } else {
      const parsedConfig = YAML.parse(this.code);
      this.code = JSON.stringify(parsedConfig, null, 2);
    }
  }

}
