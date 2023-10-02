import { Component, ViewChild } from '@angular/core';

import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { ActivatedRoute, Router } from '@angular/router';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { STRINGS } from '../../../../../assets/text/strings';
import { Tags } from 'app/common/interfaces/orb/tag';
import { CodeEditorService } from 'app/common/services/code.editor.service';
const CONFIG = {
  TAPS: 'TAPS',
  BACKEND: 'BACKEND',
  INPUTS: 'INPUTS',
  HANDLERS: 'HANDLERS',
  AGENT_POLICY: 'AGENT_POLICY',
};

@Component({
  selector: 'ngx-agent-policy-add-component',
  templateUrl: './agent.policy.add.component.html',
  styleUrls: ['./agent.policy.add.component.scss'],
})
export class AgentPolicyAddComponent {
  strings = { stepper: STRINGS.stepper };

  // #forms
  // agent policy general information - name, desc, backend
  detailsFG: FormGroup;

  // #key inputs holders
  // selected backend object
  backend: { [propName: string]: any };

  // #services responses
  // hold info retrieved
  availableBackends: {
    [propName: string]: {
      backend: string,
      description: string,
    },
  };

  agentPolicy: AgentPolicy;

  agentPolicyID: string;

  reviewPolicyConfig: boolean;

  editorVisible = true;

  errorConfigMessage: string;

  @ViewChild('editorComponentYaml')
  editorYaml;

  @ViewChild('editorComponentJson')
  editorJson;

  isEdit: boolean;

  editorOptions = {
    theme: 'vs-dark',
    language: 'yaml',
    automaticLayout: true,
    glyphMargin: false,
    folding: true,
    // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
  };
  editorOptionsJson = {
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
    readOnly: false,
    scrollBeyondLastLine: false,
    // Undocumented see https://github.com/Microsoft/vscode/issues/30795#issuecomment-410998882
    lineDecorationsWidth: 0,
    lineNumbersMinChars: 0,
  };

  codeyaml = `handlers:
  modules:
    default_dns:
      type: dns
    default_net:
      type: net
input:
  input_type: pcap
  tap: default_pcap
kind: collection`;

  codejson = 
  `{
  "handlers": {
    "modules": {
      "default_dns": {
        "type": "dns"
      },
      "default_net": {
        "type": "net"
      }
    }
  },
  "input": {
    "input_type": "pcap",
    "tap": "default_pcap"
  },
  "kind": "collection"
}
    `;

  // is config specified wizard mode or in YAML or JSON
  isJsonMode = true;

  // format definition
  format = 'yaml';

  // #load controls
  isLoading = Object.entries(CONFIG)
    .reduce((acc, [value]) => {
      acc[value] = false;
      return acc;
    }, {}) as { [propName: string]: boolean };

  selectedTags: Tags;

  uploadIconKey = 'upload-outline'

  isRequesting: boolean;  

  constructor(
    private agentPoliciesService: AgentPoliciesService,
    private notificationsService: NotificationsService,
    private router: Router,
    private route: ActivatedRoute,
    private _formBuilder: FormBuilder,
    private editor: CodeEditorService,
  ) {
    this.reviewPolicyConfig = false;
    this.isRequesting = false;
    this.agentPolicyID = this.route.snapshot.paramMap.get('id');
    this.agentPolicy = this.newAgent();
    this.isEdit = !!this.agentPolicyID;
    this.errorConfigMessage = '';

    this.readyForms();

    Promise.all([
      this.isEdit ? this.retrieveAgentPolicy() : Promise.resolve(),
      this.getBackendsList(),
    ]).catch(reason => console.warn(`Couldn't fetch data. Reason: ${reason}`))
      .then(() => this.updateForms())
      .catch((reason) => console.warn(`Couldn't fetch ${this.agentPolicy?.backend} data. Reason: ${reason}`));
  }
  ngOnInit(): void {
    this.selectedTags = this.agentPolicy?.tags || {};
  }
  resizeComponents() {
    const timeoutId = setTimeout(() => {
      window.dispatchEvent(new Event('resize'));
      clearTimeout(timeoutId);
    }, 50);
    !!this.editorJson?.layout && this.editorJson.layout();
    !!this.editorYaml?.layout && this.editorYaml.layout();
  }

  newAgent() {
    return {
      name: '',
      description: '',
      backend: 'pktvisor',
      tags: {},
      version: 1,
      policy: {
        kind: 'collection',
        input: {
          config: {},
          tap: '',
          input_type: '',
        },
        handlers: {
          modules: {},
        },
      },
    } as AgentPolicy;
  }

  retrieveAgentPolicy() {
    return new Promise(resolve => {
      this.agentPoliciesService.getAgentPolicyById(this.agentPolicyID).subscribe(policy => {
        this.agentPolicy = policy;
        this.isLoading[CONFIG.AGENT_POLICY] = false;
        resolve(policy);
      });
    });
  }

  isLoadComplete() {
    return !Object.values(this.isLoading).reduce((prev, curr) => prev || curr, false);
  }

  readyForms() {
    this.detailsFG = this._formBuilder.group({
      name: [name, [
          Validators.required,
          Validators.pattern('^[a-zA-Z_][a-zA-Z0-9_-]*$'),
          Validators.maxLength(64),
      ]],
      description: [[
        Validators.maxLength(64),
      ]],
      backend: [[Validators.required]],
    });
  }

  updateForms() {
    const {
      name,
      description,
      backend,
      format,
      policy_data,
      policy: {
        handlers,
      },
    } = this.agentPolicy;

    const wizard = format !== this.format;

    if (policy_data) {
      this.isJsonMode = false;
      this.codeyaml = policy_data;
    }

    this.detailsFG.patchValue({ name, description, backend });

    if (wizard) {
      this.onBackendSelected(backend).catch(reason => console.warn(`${reason}`));
    }
  }

  getBackendsList() {
    return new Promise((resolve) => {
      this.isLoading[CONFIG.BACKEND] = true;
      this.agentPoliciesService.getAvailableBackends().subscribe(backends => {
        this.availableBackends = !!backends && backends.reduce((acc, curr) => {
          acc[curr.backend] = curr;
          return acc;
        }, {});

        this.isLoading[CONFIG.BACKEND] = false;

        resolve(backends);
      });
    });
  }

  onBackendSelected(selectedBackend) {
    return new Promise(() => {
      this.backend = this.availableBackends[selectedBackend];
      this.backend['config'] = {};
    });
  }

  goBack() {
    this.router.navigateByUrl('/pages/datasets/policies');
  }

  viewPolicy(id) {
    this.router.navigateByUrl(`/pages/datasets/policies/view/${id}`);
  }
  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    const reader: FileReader = new FileReader();
  
    reader.onload = (e: any) => {
    const fileContent = e.target.result;
      if (this.isJsonMode) {
        this.codejson = fileContent;
      } else {
        this.codeyaml = fileContent;
      }
    };
  
    reader.readAsText(file);
  }
  onSubmit() {
    this.isRequesting = true;
    let payload = {};
    if (this.isJsonMode) {
      const policy = JSON.parse(this.codejson);
      payload = {
        name: this.detailsFG.controls.name.value,
        description: this.detailsFG.controls.description.value,
        backend: this.detailsFG.controls.backend.value,
        policy: policy,
        version: !!this.isEdit && !!this.agentPolicy.version && this.agentPolicy.version || 1,
        tags: this.selectedTags,
      }
    }
    else {
      payload = {
        name: this.detailsFG.controls.name.value,
        description: this.detailsFG.controls.description.value,
        backend: this.detailsFG.controls.backend.value,
        format: this.format,
        policy_data: this.codeyaml,
        version: !!this.isEdit && !!this.agentPolicy.version && this.agentPolicy.version || 1,
        tags: this.selectedTags,
      };
    }
    this.submit(payload);
  }

  submit(payload) {
    this.agentPoliciesService.addAgentPolicy(payload).subscribe(
      (next) => {
        this.notificationsService.success('Agent Policy successfully created', '');
        this.viewPolicy(next.id);
      },
      (error) => {
        this.notificationsService.error(
          'Failed to create Agent Policy',
          `Error: ${error.status} - ${error.statusText} - ${error.error.error}`,
        );
        this.isRequesting = false;
      },
    );   
  }
  canCreate() {
    if (this.isJsonMode) {
      if (this.editor.isJson(this.codejson)) {
        this.errorConfigMessage = '';
        return true;
      }
      else {
        this.errorConfigMessage = 'Invalid JSON configuration, check sintaxe errors';
        return false;
      }
    } else {
      if (this.editor.isYaml(this.codeyaml) && !this.editor.isJson(this.codeyaml)) {
        this.errorConfigMessage = '';
        return true;
      }
      else {
        this.errorConfigMessage = 'Invalid YAML configuration, check sintaxe errors';
        return false;
      }
    }
  }
  refreshEditor() {
    this.editorVisible = false; setTimeout(() => { this.editorVisible = true; }, 0); 
  }
  
}
