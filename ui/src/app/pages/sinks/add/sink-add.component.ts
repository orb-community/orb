import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkConfigComponent } from 'app/shared/components/orb/sink/sink-config/sink-config.component';
import { SinkDetailsComponent } from 'app/shared/components/orb/sink/sink-details/sink-details.component';
import { STRINGS } from 'assets/text/strings';
import * as YAML from 'yaml';
import { CodeEditorService } from 'app/common/services/code.editor.service';
import { SinkFeature } from 'app/common/interfaces/orb/sink/sink.feature.interface';

@Component({
    selector: 'ngx-sink-add-component',
    templateUrl: './sink-add.component.html',
    styleUrls: ['./sink-add.component.scss'],
})

export class SinkAddComponent {

    @ViewChild(SinkDetailsComponent) detailsComponent: SinkDetailsComponent;

    @ViewChild(SinkConfigComponent) configComponent: SinkConfigComponent;

    strings = STRINGS;

    createMode: boolean = true;

    sinkBackend: any;

    isRequesting: boolean;

    errorConfigMessage: string = '';

    isLoading = true;

    sinkTypesList = [];

    constructor(
        private sinksService: SinksService,
        private notificationsService: NotificationsService,
        private router: Router,
        private editor: CodeEditorService,
    ) {
        this.createMode = true;
        this.isRequesting = false;
        this.errorConfigMessage = '';
        Promise.all([this.getSinkBackends()]).then((responses) => {
            const backends = responses[0];
            this.sinkTypesList = backends.map(entry => entry.backend);
            this.isLoading = false;
        });
    }
    getSinkBackends() {
        return new Promise<SinkFeature[]>(resolve => {
          this.sinksService.getSinkBackends().subscribe(backends => {
            resolve(backends);
          });
        });
    }

    canCreate() {
        const detailsValid = this.createMode
        ? this.detailsComponent?.formGroup?.status === 'VALID'
        : true;

        const configSink = this.configComponent?.code;
        let config;

        if (this.editor.isJson(configSink)) {
            config = JSON.parse(configSink);
        } else if (this.editor.isYaml(configSink)) {
            config = YAML.parse(configSink);
            this.errorConfigMessage = '';
        } else {
            this.errorConfigMessage = 'Invalid YAML configuration, check syntax errors.';
            return false;
        }

        return !this.editor.checkEmpty(config.authentication)
        && !this.editor.checkEmpty(config.exporter)
        && detailsValid
        && !this.checkString(config);
    }
    checkString(config: any): boolean {
        if config.authentication.type === 'basicauth' {
            if (typeof config.authentication.password !== 'string' || typeof config.authentication.username !== 'string') {
                    return true;
            }
        } else if (config.authentication.type === 'noauth') {
            return true
        }
        return false;
    }

    createSink() {
        this.isRequesting = true;
        const sinkDetails = this.detailsComponent.formGroup?.value;
        const tags = this.detailsComponent.selectedTags;
        const configSink = this.configComponent.code;

        const details = { ...sinkDetails };

        let payload = {};

        const config = YAML.parse(configSink);

        payload = {
            ...details,
            tags,
            config,
        } as Sink;

        this.sinksService.addSink(payload).subscribe(() => {
            this.notificationsService.success('Sink successfully created', '');
            this.goBack();
        },
        (error) => {
          this.isRequesting = false;
        });
    }

    goBack() {
        this.router.navigateByUrl('/pages/sinks');
    }

    getBackendEmit(backend: any) {
        this.sinkBackend = backend;
    }

}
