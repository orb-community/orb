import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import 'rxjs/add/observable/empty';

import { environment } from 'environments/environment';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import {
  NgxDatabalePageInfo,
  OrbPagination,
} from 'app/common/interfaces/orb/pagination.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { expand, reduce } from 'rxjs/operators';

// default filters
const defLimit: number = 100;
const defOrder: string = 'name';
const defDir = 'asc';

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
    return this.http
      .post(
        environment.datasetPoliciesUrl,
        { ...datasetItem },
        { observe: 'response' },
      )
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to create Dataset for this Policy',
          `Error: ${err.status} - ${err.statusText} - ${err.error.error}`,
        );
        return Observable.throwError(err);
      });
  }

  getDatasetById(id: string): Observable<Dataset> {
    return this.http
      .get<Dataset>(`${environment.datasetPoliciesUrl}/${id}`)
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to fetch Dataset of this Policy',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  editDataset(datasetItem: Dataset): any {
    return this.http
      .put(`${environment.datasetPoliciesUrl}/${datasetItem.id}`, datasetItem)
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to edit Dataset',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  deleteDataset(id: string) {
    return this.http
      .delete(`${environment.datasetPoliciesUrl}/${id}`)
      .map((resp) => {
        return resp;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to Delete Dataset Policies',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }

  getAllDatasets() {
    this.clean();
    const pageInfo = DatasetPoliciesService.getDefaultPagination();

    return this.getDatasetPolicies(pageInfo).pipe(
      expand((data) => {
        return data.next
          ? this.getDatasetPolicies(data.next)
          : Observable.empty();
      }),
      reduce<OrbPagination<Dataset>>((acc, value) => {
        acc.data = value.data;
        acc.offset = 0;
        return acc;
      }, this.cache),
    );
  }

  getDatasetPolicies(pageInfo: NgxDatabalePageInfo, isFilter = false) {
    let limit = pageInfo?.limit || this.cache.limit;
    let order = pageInfo?.order || this.cache.order;
    let dir = pageInfo?.dir || this.cache.dir;
    let offset = pageInfo?.offset || 0;
    let doClean = false;
    let params = new HttpParams();

    if (isFilter) {
      if (pageInfo?.name) {
        params = params.set('name', pageInfo.name);
        // is filter different than last filter?
        doClean =
          !this.paginationCache?.name ||
          this.paginationCache?.name !== pageInfo.name;
      }
      // was filtered, no longer
    } else if (this.paginationCache?.isFilter === true) {
      doClean = true;
    }

    if (
      pageInfo.order !== this.cache.order ||
      pageInfo.dir !== this.cache.dir
    ) {
      doClean = true;
    }

    if (doClean) {
      this.clean();
      offset = 0;
      limit = this.cache.limit = pageInfo.limit;
      dir = pageInfo.dir;
      order = pageInfo.order;
    }

    if (this.paginationCache[offset]) {
      return of(this.cache);
    }
    params = params
      .set('offset', offset.toString())
      .set('limit', limit.toString())
      .set('order', order)
      .set('dir', dir);

    return this.http
      .get(environment.datasetPoliciesUrl, { params })
      .map((resp: any) => {
        this.paginationCache[pageInfo?.offset / pageInfo?.limit || 0] = true;

        // This is the position to insert the new data
        const start = pageInfo?.offset;

        const newData = [...this.cache.data];

        newData.splice(start, resp.limit, ...resp.datasets);

        this.cache = {
          ...this.cache,
          next: resp.offset + resp.limit < resp.total && {
            limit: resp.limit,
            offset: (
              parseInt(resp.offset, 10) + parseInt(resp.limit, 10)
            ).toString(),
            order: 'name',
            dir: 'desc',
          },
          limit: resp.limit,
          offset: resp.offset,
          dir: resp.direction,
          order: resp.order,
          total: resp.total,
          data: newData,
          name: pageInfo?.name,
        };

        return this.cache;
      })
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Datasets of Policy',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }
}
