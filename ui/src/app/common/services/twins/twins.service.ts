import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { Twin, PageFilters } from 'app/common/interfaces/mainflux.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Injectable()
export class TwinsService {

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) { }

  getTwins(filters: PageFilters) {
    const offset = filters.offset || 0;
    const limit = filters.limit || 10;

    const params = new HttpParams()
      .set('offset', offset.toString())
      .set('limit', limit.toString());

    return this.http.get(environment.twinsUrl, { params })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get twins',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  addTwin(twin: Twin) {
    return this.http.post(environment.twinsUrl, twin, { observe: 'response' })
      .map(
        resp => {
          this.notificationsService.success('Twin successfully created', '');
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create twin',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  getTwin(twinID: string) {
    return this.http.get(environment.twinsUrl + '/' + twinID)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch twin',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  deleteTwin(twinID: string) {
    return this.http.delete(environment.twinsUrl + '/' + twinID)
      .map(
        resp => {
          this.notificationsService.success('Twin successfully deleted', '');
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete Twin',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  editTwin(twin: Twin) {
    return this.http.put(environment.twinsUrl + '/' + twin.id, twin)
      .map(
        resp => {
          this.notificationsService.success('Twin successfully edited', '');
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Twin',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  listStates(twinID: string, offset?: number, limit?: number) {
    offset = offset || 0;
    limit = limit || 10;

    const params = new HttpParams()
      .set('offset', offset.toString())
      .set('limit', limit.toString());

    return this.http.get(environment.statesUrl + '/' + twinID, { params })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get states',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }
}
