import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { ThingsService } from 'app/common/services/things/things.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

import { MsgFilters, ReaderUrl } from 'app/common/interfaces/mainflux.interface';

const defLimit: number = 100;

@Injectable()
export class MessagesService {


  constructor(
    private http: HttpClient,
    private thingsService: ThingsService,
    private notificationsService: NotificationsService,
  ) { }

  getMessages(channel: string, thingKey: string, filters: MsgFilters, readerUrl?: ReaderUrl) {
    filters.offset = filters.offset || 0;
    filters.limit = filters.limit || defLimit;

    const headers = new HttpHeaders({
      'Authorization': thingKey,
    });

    const prefix  = readerUrl ? readerUrl.prefix : environment.readerPrefix;
    const suffix  = readerUrl ? readerUrl.suffix : environment.readerSuffix;

    let url = `${environment.readerUrl}/${prefix}/${channel}/${suffix}?`;

    Object.keys(filters).forEach(key => {
      url = filters[key] ? url += `&${key}=${filters[key]}` : url;
    });

    return this.http.get(url, { headers: headers })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to read Messages',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  sendMessage(channel: string, key: string, msg: string, subtopic?: string) {
    const headers = new HttpHeaders({
      'Authorization': key,
    });

    let url = `${environment.httpAdapterUrl}/${environment.readerPrefix}/${channel}/${environment.readerSuffix}`;
    url = subtopic ? url += `/${encodeURIComponent(subtopic)}` : url;

    return this.http.post(url, msg, { headers: headers })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to send Message',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  sendLocationMock(chanID: string, thingID: string) {
    const lon = 44.7 + 0.5 * Math.random();
    const lat = 20.4 + 0.5 * Math.random();

    const message = `[{"bn":"location-", "n":"lon", "v":${lon}}, {"n":"lat", "v":${lat}}]`;

    this.thingsService.getThing(thingID).subscribe(
      (resp: any) => {
        this.sendMessage(chanID, resp.key, message).subscribe();
      },
    );
  }
}
