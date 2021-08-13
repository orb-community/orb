/**
 * @license
 * Copyright Akveo. All Rights Reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 */
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { NgModule } from '@angular/core';
import { HttpClientModule } from '@angular/common/http';
import { CoreModule } from './@core/core.module';
import { ThemeModule } from './@theme/theme.module';
import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import {
  NbAlertModule,
  NbButtonModule,
  NbCardModule,
  NbChatModule,
  NbCheckboxModule,
  NbDatepickerModule,
  NbDialogModule,
  NbIconModule,
  NbInputModule,
  NbLayoutModule,
  NbMenuModule,
  NbSidebarModule,
  NbToastrModule,
  NbWindowModule,
} from '@nebular/theme';

// MFx- Foorm dependency
import { FormsModule } from '@angular/forms';
// Mfx - MQTT dependencies for Gateways page
import { IMqttServiceOptions, MqttModule, MqttService } from 'ngx-mqtt';
import { environment } from 'environments/environment';
// Mfx - Auth and Profile pages
import { BreadcrumbModule } from 'xng-breadcrumb';
import { NgxDatatableModule } from '@swimlane/ngx-datatable';

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  connectOnCreate: false,
  url: environment.mqttWsUrl,
};

@NgModule({
  declarations: [
    AppComponent,
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    AppRoutingModule,

    ThemeModule.forRoot(),

    NbSidebarModule.forRoot(),
    NbMenuModule.forRoot(),
    NbDatepickerModule.forRoot(),
    NbDialogModule.forRoot(),
    NbWindowModule.forRoot(),
    NbToastrModule.forRoot(),
    NbChatModule.forRoot({
      messageGoogleMapKey: 'AIzaSyA_wNuCzia92MAmdLRzmqitRGvCF7wCZPY',
    }),
    CoreModule.forRoot(),
    // Mfx dependencies
    MqttModule.forRoot(MQTT_SERVICE_OPTIONS),
    FormsModule,
    NbInputModule,
    NbCardModule,
    NbIconModule,
    NbButtonModule,

    // 3rd party
    BreadcrumbModule,
    NgxDatatableModule,
    NbAlertModule,
    NbCheckboxModule,
  ],
  bootstrap: [AppComponent],
  // Mfx dependencies
  providers: [MqttService],
})
export class AppModule {
}
