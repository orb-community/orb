import { Component, OnDestroy, OnInit } from '@angular/core';
import { NbDialogRef } from '@nebular/theme';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { DynamicFormConfig } from 'app/common/interfaces/orb/dynamic.form.interface';
import { PolicyHandler } from 'app/common/interfaces/orb/policy/policy.handler.interface';
import { PolicyBackend } from 'app/common/interfaces/orb/policy/policy.backend.interface';

type ConfigHandler = PolicyHandler & {
  config?: DynamicFormConfig,
  filter?: DynamicFormConfig,
  metrics?: DynamicFormConfig,
  metrics_groups?: DynamicFormConfig,
};

interface HandlerCollection {
  [propName: string]: ConfigHandler;
}

interface ModuleCollection {
  [propName: string]: PolicyHandler;
}

@Component({
  selector: 'ngx-agent-policy-details-component',
  templateUrl: './handler.policy.add.component.html',
  styleUrls: ['./handler.policy.add.component.scss'],
})
export class HandlerPolicyAddComponent implements OnInit, OnDestroy {
  // handlers
  handlerSelectorFG: FormGroup;

  // backend - selected by user on agent policy creation
  backend: PolicyBackend;

  // holds all handlers added by user
  modules: ModuleCollection = {};

  // handler key
  selectedHandler: string;

  // handler dyn configs to render
  dynConfigList = ['config', 'filter'];

  handlerProps: ConfigHandler;

  isLoading: boolean;

  subscription: Subscription;

  availableHandlers: HandlerCollection = {};

  constructor(
    protected dialogRef: NbDialogRef<HandlerPolicyAddComponent>,
    protected route: ActivatedRoute,
    protected router: Router,
    protected _formBuilder: FormBuilder,
    protected agentPoliciesService: AgentPoliciesService,
  ) {
    this.isLoading = true;

    this.backend = this.router.getCurrentNavigation().extras.state?.backend as PolicyBackend || null;
    this.modules = this.router.getCurrentNavigation().extras.state?.modules as ModuleCollection || null;

    this.subscription = this.getHandlers()
      .subscribe(handlers => {
        this.availableHandlers = handlers;
        this.isLoading = false;
      });
  }

  ngOnInit() {
    this.readyForms();
  }

  ngOnDestroy() {
    this.subscription.unsubscribe();
  }

  onClose() {
    this.dialogRef.close(false);
  }

  getHandlers() {
    return this.agentPoliciesService.getBackendConfig([this.backend.backend, 'handlers'])
      .map(handlers => Object.entries<[string, { HandlerCollection }]>(handlers)
          .reduce((acc, [key, value]) => {
            const latest = Object.entries(value as [string, HandlerCollection])
              .sort(([a], [b]) => a < b ? -1 : a > b ? 1 : 0)
              .map(([version, content]) => ({ content, version }))
              .pop();
            acc[key] = latest;
            return acc;
          }, {})
        || {});
  }

  readyForms(): void {
    this.handlerSelectorFG = this._formBuilder.group({
      'selected_handler': [null, [Validators.required]],
      'label': [null, [Validators.required]],
      'type': [null],
      'config': [null],
      'filter': [null],
    });
  }

  createDynamicControls(config) {
    const controlReducer = (previous, [key, __]) => {
      previous[key] = ['', [Validators.required]];
      return previous;
    };

    const dynamicControls = Object.entries(config || {}).reduce(controlReducer, {});

    const dynamicFormGroup = Object.keys(dynamicControls).length > 0 ? this._formBuilder.group(dynamicControls) : null;

    return dynamicFormGroup;
  }

  onHandlerSelected(selectedHandler) {

    const { config, filter, type } = this.handlerProps = this.availableHandlers[selectedHandler];

    const suggestName = this.getSeriesHandlerName(selectedHandler);

    this.handlerSelectorFG.patchValue({
      label: suggestName,
      type,
      config: this.createDynamicControls(config),
      filter: this.createDynamicControls(filter),
    });
  }

  getSeriesHandlerName(handlerType) {
    const count = 1 + Object.values(this.modules || {}).filter(({ type }) => type === handlerType).length;
    return `handler_${ handlerType }_${ count }`;
  }

  checkValidName() {
    const { value } = this.handlerSelectorFG.controls.label;
    const hasTagForKey = Object.keys(this.modules).find(key => key === value);
    return value && value !== '' && !hasTagForKey;
  }

  onSaveHandler() {
    const configForm = this.handlerSelectorFG.get('config') as FormGroup;
    const filterForm = this.handlerSelectorFG.get('filter') as FormGroup;
    const { label, type } = this.handlerSelectorFG.value;
    let config, filter;

    if (configForm !== null) {
      config = Object.entries(configForm.controls)
        .reduce((acc, [key, control]) => {
          if (control.value) acc[key] = control.value;
          return acc;
        }, {});
    }

    if (filterForm !== null) {
      filter = Object.entries(filterForm.controls)
        .reduce((acc, [key, control]) => {
          if (control.value) acc[key] = control.value;
          return acc;
        }, {});
    }

    this.dialogRef.close({
      label,
      type,
      config,
      filter,
    });
  }
}
