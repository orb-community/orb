import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { Sink } from 'app/common/interfaces/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { PageFilters } from 'app/common/interfaces/mainflux.interface';

// default filters
const defLimit: number = 20;
const defOrder: string = 'id';
const defDir: string = 'desc';

@Injectable()
export class SinksService {
  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
  }

  addSink(sinkItem: Sink) {
    return this.http.post(environment.sinksUrl,
      sinkItem,
      {observe: 'response'})
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Sink',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
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
          return Observable.throwError(err);
        },
      );
  }

  getSinks(filters: PageFilters) {
    filters.offset = filters.offset || 0;
    filters.limit = filters.limit || defLimit;
    filters.order = filters.order || defOrder;
    filters.dir = filters.dir || defDir;

    let params = new HttpParams()
      .set('offset', filters.offset.toString())
      .set('limit', filters.limit.toString())
      .set('order', filters.order)
      .set('dir', 'asc');

    if (filters.name) {
      params = params.append('name', filters.name);
    }

    return this.http.get(environment.sinksUrl, {params})
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get sinks',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  editSink(sinkItem: Sink): any {
    return this.http.put(`${environment.sinksUrl}/${sinkItem.id}`, sinkItem)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Sink',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }

  deleteSink(sinkId: string) {
    return this.http.delete(`${environment.sinksUrl}/${sinkId}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete Sink',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throwError(err);
        },
      );
  }
}
