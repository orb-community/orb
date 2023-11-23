import { Component, EventEmitter, Input, OnChanges, OnInit, Output, SimpleChanges, ViewChild } from '@angular/core';
import { FormBuilder, FormControl, Validators } from '@angular/forms';
import { Sink, SinkBackends } from 'app/common/interfaces/orb/sink.interface';
import * as YAML from 'yaml';
import IStandaloneEditorConstructionOptions = monaco.editor.IStandaloneEditorConstructionOptions;
import { OrbService } from 'app/common/services/orb.service';
import { SINK_OTLP_CONFIG_YAML, SINK_PROMETHEUS_CONFIG_YAML } from 'app/shared/configurations/configurations_schemas';

@Component({
  selector: 'ngx-sink-config',
  templateUrl: './sink-config.component.html',
  styleUrls: ['./sink-config.component.scss'],
})
export class SinkConfigComponent implements OnInit, OnChanges {

  @Input()
  sink: Sink;

  @Input()
  editMode: boolean = false;

  @Input()
  createMode: boolean = false;

  @Input()
  sinkBackend: string;

  @Output()
  editModeChange: EventEmitter<boolean>;

  @Input()
  detailsEditMode: boolean;

  @Input()
  errorConfigMessage: string;

  @ViewChild('editorComponent')
  editor;

  selectBackendWarning = false;
  warningMessageTop = 0;
  warningMessageLeft = 0;

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
    this.errorConfigMessage = '';
  }

  ngOnInit(): void {
    if (this.createMode) {
      this.updateForm();
      this.code = SINK_OTLP_CONFIG_YAML;
    } else {
      const parsedCode = YAML.parse(JSON.stringify(this.sink.config));
      this.code = YAML.stringify(parsedCode);
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
    this.code = this.sinkBackend === SinkBackends.prometheus ? SINK_PROMETHEUS_CONFIG_YAML : SINK_OTLP_CONFIG_YAML;
    if (sinkBackend.currentValue) {
      this.editorOptionsYaml = { ...this.editorOptionsYaml, readOnly: false };
    }
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
  onEditorClick(event: any) {
    if (this.createMode && this.sinkBackend === undefined && !this.selectBackendWarning) {
      this.selectBackendWarning = true;
      setTimeout(() => {
        this.selectBackendWarning = false;
      }, 2000);
    }
    this.warningMessageTop = event.target.clientY;
    this.warningMessageLeft = event.target.clientX;
  }
}
