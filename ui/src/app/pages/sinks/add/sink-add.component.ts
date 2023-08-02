import { Component, ViewChild } from '@angular/core';
import { Router } from '@angular/router';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import { SinkConfigComponent } from 'app/shared/components/orb/sink/sink-config/sink-config.component';
import { SinkDetailsComponent } from 'app/shared/components/orb/sink/sink-details/sink-details.component';
import { STRINGS } from 'assets/text/strings';

@Component({
    selector: 'ngx-sink-add-component',
    templateUrl: './sink-add.component.html',
    styleUrls: ['./sink-add.component.scss'],
})

export class SinkAddComponent {

    @ViewChild(SinkDetailsComponent) detailsComponent: SinkDetailsComponent;

    @ViewChild(SinkConfigComponent) configComponent: SinkConfigComponent;

    strings = STRINGS;

    createMode: boolean;

    sinkBackend: any;

    constructor(
        private sinksService: SinksService,
        private notificationsService: NotificationsService,
        private router: Router,
    ) {
        this.createMode = true;
    }

    canCreate() { 
        const detailsValid = this.createMode
        ? this.detailsComponent?.formGroup?.status === 'VALID'
        : true;
        return detailsValid;
    }

    createSink() {
        const sinkDetails = this.detailsComponent.formGroup?.value;
        const tags = this.detailsComponent.selectedTags;
        const configSink = this.configComponent.code;

        const details = { ...sinkDetails };
        
        let payload = {};

        if (this.isJson(configSink)) {
            const config = JSON.parse(configSink);
            payload = {
                ...details,
                tags,
                config,
            } as Sink;
        }
        else {
            payload = {
                ...details,
                tags,
                format: 'yaml',
                config_data: configSink,
            } as Sink;
        }

        this.sinksService.addSink(payload).subscribe(() => {
            this.notificationsService.success('Sink successfully created', '');
            this.goBack();
        });
    }
    isJson(str: string) {
        try {
            JSON.parse(str);
            return true;
        } catch {
            return false;
        }
    }

    goBack() {
        this.router.navigateByUrl('/pages/sinks');
    }

    getBackendEmit(backend: any) {
        this.sinkBackend = backend;
    }
}
