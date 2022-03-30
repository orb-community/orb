import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { NgxDatabalePageInfo, OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { delay, expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 100;
const defOrder: string = 'name';
const defDir = 'desc';

@Injectable()
export class DatasetPoliciesService {
  paginationCache: any = {};

  cache: OrbPagination<Dataset>;

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {
    this.clean();
  }

  public static getDefaultPagination(): OrbPagination<Dataset> {
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

  addDataset(datasetItem: Dataset) {
    return this.http.post(environment.datasetPoliciesUrl,
        { ...datasetItem },
        { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create Dataset Policy',
            `Error: ${ err.status } - ${ err.statusText } - ${ err.error.error }`);
          return Observable.throwError(err);
        },
      );
  }

  getDatasetById(id: string): Observable<Dataset> {
    return this.http.get<Dataset>(`${ environment.datasetPoliciesUrl }/${ id }`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Dataset Policy',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  editDataset(datasetItem: Dataset): any {
    return this.http.put(`${ environment.datasetPoliciesUrl }/${ datasetItem.id }`, datasetItem)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit Dataset Policy',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  deleteDataset(id: string) {
    return this.http.delete(`${ environment.datasetPoliciesUrl }/${ id }`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Delete Dataset Policies',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }

  getAllDatasets() {
    const pageInfo = DatasetPoliciesService.getDefaultPagination();
    pageInfo.limit = 100;

    return this.getDatasetPolicies(pageInfo)
      .pipe(
        expand(data => {
          return data.next ? this.getDatasetPolicies(data.next) : Observable.empty();
        }),
        delay(250),
        reduce<OrbPagination<Dataset>>((acc, value) => {
          acc.data.splice(value.offset, value.limit, ...value.data);
          acc.offset = 0;
          return acc;
        }, this.cache),
      );
  }

  getDatasetPolicies(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    const offset = pageInfo?.offset || this.cache.offset;
    const limit = pageInfo?.limit || this.cache.limit;
    let params = new HttpParams()
      .set('offset', (offset * limit).toString())
      .set('limit', limit.toString())
      .set('order', this.cache.order)
      .set('dir', this.cache.dir);

    if (isFilter) {
      if (pageInfo?.name) {
        params = params.append('name', pageInfo.name);
      }
      if (pageInfo?.tags) {
        params.append('tags', JSON.stringify(pageInfo.tags));
      }
      this.paginationCache[offset] = false;
    }

    if (this.paginationCache[pageInfo?.offset]) {
      return of(this.cache);
    }

    return this.http.get(environment.datasetPoliciesUrl, { params })
      .map((resp: any) => {
          this.paginationCache[pageInfo?.offset || 0] = true;
          // This is the position to insert the new data
          const start = resp.offset;
          const newData = [...this.cache.data];
          // TODO find out the field name for data in response json
          newData.splice(start, resp.limit, ...resp.datasets);
          this.cache = {
            ...this.cache,
            offset: Math.floor(resp.offset / resp.limit),
            total: resp.total,
            data: newData,
          };
          if (pageInfo?.name) this.cache.name = pageInfo.name;
          if (pageInfo?.tags) this.cache.tags = pageInfo.tags;
          return this.cache;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to get Dataset Policies',
            `Error: ${ err.status } - ${ err.statusText }`);
          return Observable.throwError(err);
        },
      );
  }
}
