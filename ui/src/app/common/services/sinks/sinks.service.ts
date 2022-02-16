import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';

// default filters
const defLimit: number = 10;
const defOrder: string = 'name';
const defDir = 'desc';

@Injectable()
export class SinksService {
  paginationCache: any = {};

  cache: OrbPagination<Sink>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
  }

  public static getDefaultPagination(): OrbPagination<Sink> {
    return {
      limit: defLimit,
      order: defOrder,
      dir: defDir,
      offset: 0,
      total: 0,
      data: null,
    };
  }

  clean() {
    this.cache = {
      limit: defLimit,
      offset: 0,
      order: defOrder,
      total: 0,
      dir: defDir,
      data: [],
    };
    this.paginationCache = {};
  }

  addSink(sinkItem: Sink) {
    return this.http.post(environment.sinksUrl,
        sinkItem,
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Sink',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  getSinkById(sinkId: string): any {
    return this.http.get(`${ environment.sinksUrl }/${ sinkId }`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Sink',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getSinkBackends() {
    return this.http.get(environment.sinkBackendsUrl)
      .map((resp: any) => {
        return resp.backends;
      }).catch(
        err => {
          this.notificationsService.error('Failed to get Sink Backends',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getSinks(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    let limit = pageInfo?.limit || this.cache.limit;
    let order = pageInfo?.order || this.cache.order;
    let dir = pageInfo?.dir || this.cache.dir;
    let offset = pageInfo?.offset || this.cache.offset;
    let doClean = false;
    let params = new HttpParams();

    if (isFilter) {
      if (pageInfo?.name) {
        params = params.set('name', pageInfo.name);
        // is filter different than last filter?
        doClean = !this.paginationCache?.name || this.paginationCache?.name !== pageInfo.name;
      }
      // was filtered, no longer
    } else if (this.paginationCache?.isFilter === true) {
      doClean = true;
    }

    if (pageInfo.order !== this.cache.order
      || pageInfo.dir !== this.cache.dir) {
      doClean = true;
    }

    if (doClean) {
      this.clean();
      offset = this.cache.offset;
      limit = this.cache.limit = pageInfo.limit;
      dir = pageInfo.dir;
      order = pageInfo.order;
    }

    if (this.paginationCache[offset]) {
      return of(this.cache);
    }

    params = params
      .set('offset', (offset * limit).toString())
      .set('limit', limit.toString())
      .set('order', order)
      .set('dir', dir);

    return this.http.get(environment.sinksUrl, { params })
      .map(
        (resp: any) => {
          this.paginationCache[pageInfo?.offset / pageInfo?.limit || 0] = true;
          // This is the position to insert the new data
          const start = pageInfo?.offset * pageInfo?.limit || 0;
          const newData = [...this.cache.data];
          newData.splice(start, resp.limit, ...resp.sinks);
          this.cache = {
            ...this.cache,
            offset: resp.offset,
            dir: resp.dir,
            order: resp.order,
            total: resp.total,
            data: newData,
            name: pageInfo?.name,
          };

          return this.cache;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Sinks',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  editSink(sinkItem: Sink): any {
    return this.http.put(`${ environment.sinksUrl }/${ sinkItem.id }`, sinkItem)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Sink',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  deleteSink(sinkId: string) {
    return this.http.delete(`${ environment.sinksUrl }/${ sinkId }`)
      .map(
        resp => {
          this.cache.data.splice(this.cache.data.map(s => s.id).indexOf(sinkId), 1);
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete Sink',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }
}
