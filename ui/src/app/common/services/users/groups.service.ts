import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { environment } from 'environments/environment';
import { UserGroup } from 'app/common/interfaces/mainflux.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

const defLimit: number = 20;

@Injectable()
export class UserGroupsService {

  constructor(
    private http: HttpClient,
    private notificationsService: NotificationsService,
  ) { }

  addGroup(org: UserGroup) {
    return this.http.post(environment.groupsUrl, org, { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create User Group',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  editGroup(org: UserGroup) {
    return this.http.patch(`${environment.groupsUrl}/${org.id}`, org)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit User Group',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  deleteGroup(orgID: string) {
    return this.http.delete(`${environment.groupsUrl}/${orgID}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to delete User Group',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  assignUser(groupID: string, userIDs: string[]): any {
    const assignReq = {
      members: userIDs,
      type: 'user',
    };
    return this.http.post(`${environment.groupsUrl}/${groupID}/members`, assignReq)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Assing User to Group',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  unassignUser(groupID: string, userIDs: string[]): any {
    const unassignReq = {
      members: userIDs,
    };
    return this.http.request('delete', `${environment.groupsUrl}/${groupID}/members`, {body: unassignReq})
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to Unassing User from Group',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  getGroup(groupID: string): any {
    return this.http.get(`${environment.groupsUrl}/${groupID}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch User Group',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  getGroups(offset?: number, limit?: number, name?: string): any {
    offset = offset || 0;
    limit = limit || defLimit;

    let params = new HttpParams()
      .set('offset', offset.toString())
      .set('limit', limit.toString());

    if (name) {
      params = params.append('name', name);
    }

    return this.http.get(environment.groupsUrl, { params })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch User Groups',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  getMembers(groupID?: string): any {
    return this.http.get(`${environment.groupsUrl}/${groupID}/members`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Group members',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }
}
