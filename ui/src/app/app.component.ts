/**
 * @license
 * Copyright Akveo. All Rights Reserved.
 * Licensed under the MIT License. See License.txt in the project root for license information.
 */
import { Component, OnInit } from '@angular/core';
import { GoogleTagManagerService } from 'angular-google-tag-manager';
import { NavigationEnd, Router } from '@angular/router';

@Component({
  selector: 'ngx-app',
  template: '<router-outlet></router-outlet>',
})
export class AppComponent implements OnInit {

  constructor(private router: Router,
              private gtmService: GoogleTagManagerService) {
  }

  ngOnInit() {
    // I've commented about this on issue [#209](https://github.com/ns1labs/orb/issues/209)
    // this.analytics.trackPageViews();
    /**
     * Track pages with router events
     * Lets test this and think of other events to track down the road
     * angular-google-tag-manager module might still be dropped out
     * in favor of analyticsService that came bundled with mainfluxUI codebase
     */
    this.router.events.forEach(item => {
      if (item instanceof NavigationEnd) {
        const gtmTag = {
          event: 'page',
          pageName: item.url,
        };

        this.gtmService.pushTag(gtmTag);
      }
    });
  }
}
