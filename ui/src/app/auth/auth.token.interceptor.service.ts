import { Injectable, Injector } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse,
} from '@angular/common/http';
import { Observable } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import 'rxjs/add/operator/switchMap';
import { tap } from 'rxjs/operators';
import { Router } from '@angular/router';

import { NbAuthJWTToken, NbAuthService } from '@nebular/auth';

import { environment } from 'environments/environment';

@Injectable()
export class TokenInterceptor implements HttpInterceptor {

  constructor(
    private inj: Injector,
    private authService: NbAuthService,
    private router: Router,
  ) { }

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    this.authService = this.inj.get(NbAuthService);

    return this.authService.getToken().pipe(switchMap(
      (token: NbAuthJWTToken) => {
        if (token && token.getValue() &&
          !request.url.startsWith(environment.httpAdapterUrl) &&
          !request.url.startsWith(environment.readerUrl) &&
          !request.url.startsWith(environment.bootstrapUrl) &&
          !request.url.startsWith(environment.browseUrl)
        ) {
          request = request.clone({
            setHeaders: {
              'Authorization': token.getValue(),
            },
          });
        }
        return next.handle(request).pipe(tap(
          resp => {
          },
          err => {
            // Status 403 - Forbiden
            if (err instanceof HttpErrorResponse && err.status === 403 &&
              !request.url.startsWith(environment.httpAdapterUrl) &&
              !request.url.startsWith(environment.readerUrl)) {
              localStorage.removeItem('auth_app_token');
              this.router.navigateByUrl('/auth/login');
            }
          },
        ));
      },
    ));
  }
}
