/**
 * @license
 * Copyright Akveo. All Rights Reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 */
import { Component, OnInit } from '@angular/core';
import { GoogleAnalyticsService } from './common/services/analytics/google-service-analytics.service';
import { environment } from 'environments/environment';

@Component({
  selector: 'ngx-app',
  template: '<router-outlet></router-outlet>',
})
export class AppComponent implements OnInit {

  constructor(private gtagService: GoogleAnalyticsService) {
  }

  ngOnInit() {
    if (!!environment.production) {
      this.gtagService.setGtagID(environment.GTAGID);
      this.gtagService.init();
    }
  }
}
