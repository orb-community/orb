import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import 'rxjs/add/observable/empty';
import { Router } from '@angular/router';

import { environment } from 'environments/environment';
import { User } from 'app/common/interfaces/mainflux.interface';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

const defLimit: number = 20;

@Injectable()
export class UsersService {
  picture = 'assets/images/mainflux-logo.png';

  constructor(
    private http: HttpClient,
    private router: Router,
    private notificationsService: NotificationsService,
  ) { }

  addUser(user: User) {
    return this.http.post(environment.usersUrl, user, { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to create User',
            `Error: ${err.status} - ${err.statusText}`);
          return Observable.throw(err);
        },
      );
  }

  getProfile(): any {
    return this.getUser('profile');
  }

  getUser(userID: string): any {
    return this.http.get(`${environment.usersUrl}/${userID}`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.router.navigateByUrl('/auth/login');
          return Observable.empty();
        },
      );
  }

  getUsers(offset?: number, limit?: number, email?: string): any {
    offset = offset || 0;
    limit = limit || defLimit;

    let params = new HttpParams()
      .set('offset', offset.toString())
      .set('limit', limit.toString());

    if (email) {
      params = params.append('email', email);
    }

    return this.http.get(environment.usersUrl, { params })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Users',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  editUser(user: User): any {
    return this.http.put(environment.usersUrl, user)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to edit User',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  changeUserPassword(passReq: any): any {
    return this.http.patch(environment.changePassUrl, passReq, { observe: 'response' })
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to change User password',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }

  getServiceVersion() {
    return this.http.get(environment.usersVersionUrl);
  }

  getUserPicture(): any {
    return this.picture;
  }

  getMemberships(memberID?: string): any {
    return this.http.get(`${environment.membersUrl}/${memberID}/groups`)
      .map(
        resp => {
          return resp;
        },
      )
      .catch(
        err => {
          this.notificationsService.error('Failed to fetch Group memberships',
            `Error: ${err.status} - ${err.statusText}`);
            return Observable.throw(err);
        },
      );
  }
}
