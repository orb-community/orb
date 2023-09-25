import {
  HttpErrorResponse,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable, Injector } from '@angular/core';
import { Router } from '@angular/router';
import { Observable } from 'rxjs';
import 'rxjs/add/operator/switchMap';
import { switchMap, tap } from 'rxjs/operators';

import { NbAuthJWTToken, NbAuthService } from '@nebular/auth';

import { environment } from 'environments/environment';

import { JwtHelperService } from '@auth0/angular-jwt';

const jwtHelper = new JwtHelperService();
@Injectable()
export class TokenInterceptor implements HttpInterceptor {
  constructor(
    private inj: Injector,
    private authService: NbAuthService,
    private router: Router,
  ) {}

  intercept(
    request: HttpRequest<any>,
    next: HttpHandler,
  ): Observable<HttpEvent<any>> {
    this.authService = this.inj.get(NbAuthService);

    return this.authService.getToken().pipe(
      switchMap((token: NbAuthJWTToken) => {
        if (
          token &&
          token.getValue() &&
          !request.url.startsWith(environment.httpAdapterUrl) &&
          !request.url.startsWith(environment.readerUrl) &&
          !request.url.startsWith(environment.bootstrapUrl) &&
          !request.url.startsWith(environment.browseUrl)
        ) {
          const expirationDate = jwtHelper.getTokenExpirationDate(
            token.getValue(),
          );
          if (expirationDate < new Date()) {
            this.authService.logout('email');
            this.router.navigate(['auth/login']);
          }

          request = request.clone({
            setHeaders: {
              Authorization: `Bearer ${token.getValue()}`,
            },
          });
        }
        return next.handle(request).pipe(
          tap(
            (resp) => {},
            (err) => {
              // Status 401 - Unauthorized
              if (
                err instanceof HttpErrorResponse &&
                err.status === 401 &&
                !request.url.startsWith(environment.httpAdapterUrl) &&
                !request.url.startsWith(environment.readerUrl)
              ) {
                localStorage.removeItem('auth_app_token');
                this.router.navigateByUrl('/auth/login');
              }
            },
          ),
        );
      }),
    );
  }
}
