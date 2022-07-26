import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { environment } from 'environments/environment';
import { catchError, expand, map, scan, takeWhile } from 'rxjs/operators';

@Injectable()
export class SinksService {
  paginationCache: any = {};

  cache: OrbPagination<Sink>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {}

  addSink(sinkItem: Sink) {
    return this.http
      .post(environment.sinksUrl, sinkItem, { observe: 'response' })
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to create Sink',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  getSinkById(sinkId: string): Observable<Sink> {
    return this.http.get(`${environment.sinksUrl}/${sinkId}`).pipe(
      catchError((err) => {
        this.notificationsService.error(
          'Failed to fetch Sink',
          `Error: ${err.status} - ${err.statusText}`,
        );
        err['id'] = sinkId;
        return of(err);
      }),
    );
  }

  getSinkBackends() {
    return this.http
      .get(environment.sinkBackendsUrl)
      .map((resp: any) => {
        return resp.backends;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Sink Backends',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAllSinks() {
    const page = {
      order: 'name',
      dir: 'asc',
      limit: 100,
      data: [],
      offset: 0,
    } as OrbPagination<Sink>;

    return this.getSinks(page).pipe(
      expand((data) => {
        return data.next ? this.getSinks(data.next) : Observable.empty();
      }),
      takeWhile((data) => data.next !== undefined),
      map((_page) => _page.data),
      scan((acc, v) => [...acc, ...v]),
    );
  }

  getSinks(page: OrbPagination<Sink>) {
    const params = new HttpParams()
      .set('order', page.order)
      .set('dir', page.dir)
      .set('offset', page.offset.toString())
      .set('limit', page.limit.toString());

    return this.http
      .get(environment.sinksUrl, { params })
      .pipe(
        map((resp: any) => {
          const {
            order,
            direction: dir,
            offset,
            limit,
            total,
            sinks: data,
            tags,
          } = resp;
          const next = offset + limit < total && {
            limit,
            order,
            dir,
            tags,
            offset: (parseInt(offset, 10) + parseInt(limit, 10)).toString(),
          };
          return {
            order,
            dir,
            offset,
            limit,
            total,
            data,
            next,
          } as OrbPagination<Sink>;
        }),
      )
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Sinks',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  editSink(sinkItem: Sink): any {
    return this.http
      .put(`${environment.sinksUrl}/${sinkItem.id}`, sinkItem)
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to edit Sink',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  deleteSink(sinkId: string) {
    return this.http
      .delete(`${environment.sinksUrl}/${sinkId}`)
      .map((resp) => {
        this.cache.data.splice(
          this.cache.data.map((s) => s.id).indexOf(sinkId),
          1,
        );
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to delete Sink',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }
}
