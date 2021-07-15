import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { Sink } from 'app/common/interfaces/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { PageFilters } from 'app/common/interfaces/mainflux.interface';

const defLimit: number = 20;

@Injectable()
export class SinksService {
  picture = 'assets/images/mainflux-logo.png';

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) { }

  addSink(sink: Sink) {
    return this.http.post(environment.sinksUrl, sink, { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Sink',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  getSinkById(sinkId: string): any {
    return this.http.get(`${environment.sinksUrl}/${sinkId}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
           this.notificationsService.error('Failed to fetch Sink',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  getSinks(filters: PageFilters) {
    filters.offset = filters.offset || 0;
    filters.limit = filters.limit || defLimit;

    let params = new HttpParams()
      .set('offset', filters.offset.toString())
      .set('limit', filters.limit.toString())
      .set('order', 'name')
      .set('dir', 'asc');

    if (filters.type) {
      if (filters.metadata) {
        params = params.append('metadata', `{"${filters.type}": ${filters.metadata}}`);
      } else {
        params = params.append('metadata', `{"type":"${filters.type}"}`);
      }
    }

    if (filters.name) {
      params = params.append('name', filters.name);
    }

    return this.http.get(environment.sinksUrl, { params })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get sinks',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  editSink(sinkItem: Sink): any {
    return this.http.put(environment.sinksUrl, sinkItem)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Sink',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

 deleteSink(sinkId: string) {
    return this.http.delete(`${environment.thingsUrl}/${sinkId}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete Sink',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }
}
