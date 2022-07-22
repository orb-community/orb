import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';

import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { OrbPagination } from 'app/common/interfaces/orb/pagination.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';
import { environment } from 'environments/environment';
import { expand, map, scan, takeWhile } from 'rxjs/operators';

@Injectable()
export class DatasetPoliciesService {
  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) {}

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
    const page = {
      order: 'name',
      dir: 'asc',
      limit: 100,
      data: [],
      offset: 0,
    } as OrbPagination<Dataset>;

    return this.getDatasetPolicies(page).pipe(
      expand((data) => {
        return data.next
          ? this.getDatasetPolicies(data.next)
          : Observable.empty();
      }),
      takeWhile((data) => data.next !== undefined),
      map((_page) => _page.data),
      scan((acc, v) => [...acc, ...v]),
    );
  }

  getDatasetPolicies(page: OrbPagination<Dataset>) {
    const params = new HttpParams()
      .set('order', page.order)
      .set('dir', page.dir)
      .set('offset', page.offset.toString())
      .set('limit', page.limit.toString());

    return this.http
      .get(environment.datasetPoliciesUrl, { params })
      .pipe(
        map((resp: any) => {
          const {
            order,
            direction: dir,
            offset,
            limit,
            total,
            datasets: data,
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
          } as OrbPagination<Dataset>;
        }),
      )
      .catch((err) => {
        this.notificationsService.error(
          'Failed to get Datasets',
          `Error: ${err.status} - ${err.statusText}`,
        );
        return Observable.throwError(err);
      });
  }
}
