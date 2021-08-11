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
  NbButtonModule,
  NbCardModule,
  NbChatModule,
  NbDatepickerModule,
  NbDialogModule,
  NbIconModule,
  NbInputModule,
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
import { LogoutComponent } from './pages/logout/logout.component';
import { RegisterComponent } from './pages/register/register.component';
import { ProfileComponent } from './pages/profile/profile.component';
import { BreadcrumbModule } from 'xng-breadcrumb';

export const MQTT_SERVICE_OPTIONS: IMqttServiceOptions = {
  connectOnCreate: false,
  url: environment.mqttWsUrl,
};

@NgModule({
  declarations: [
    AppComponent,
    // Mfx Componennt
    LogoutComponent,
    RegisterComponent,
    ProfileComponent,
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
  ],
  bootstrap: [AppComponent],
  // Mfx dependencies
  providers: [MqttService],
})
export class AppModule {
}
