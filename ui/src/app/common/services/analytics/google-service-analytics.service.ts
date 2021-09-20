import { Injectable } from '@angular/core';
import { NavigationEnd, Router } from '@angular/router';

declare var gtag: Function;

@Injectable({
  providedIn: 'root',
})
export class GoogleAnalyticsService {
  private GTAGID: String;

  constructor(private router: Router) {
  }

  public event(eventName: string, params: {}) {
    gtag('event', eventName, params);
  }

  public init() {
    this.listenForRouteChanges();

    try {

      const script1 = document.createElement('script');
      script1.async = true;
      script1.src = 'https://www.googletagmanager.com/gtag/js?id=' + this.GTAGID;
      document.head.appendChild(script1);

      const script2 = document.createElement('script');
      script2.innerHTML = `
        window.dataLayer = window.dataLayer || [];
        function gtag(){dataLayer.push(arguments);}
        gtag('js', new Date());
        gtag('config', '${this.GTAGID}');
      `;
      document.head.appendChild(script2);
    } catch (ex) {
      console.error('Error appending google analytics');
      console.error(ex);
    }
  }

  private listenForRouteChanges() {
    this.router.events.subscribe(event => {
      if (event instanceof NavigationEnd) {
        gtag('config', this.GTAGID, {
          'page_path': event.urlAfterRedirects,
        });
      }
    });
  }

  public setGtagID(gtagID: String) {
    this.GTAGID = gtagID;
  }
}
