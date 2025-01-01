import {HttpEvent, HttpHandlerFn, HttpRequest} from "@angular/common/http";
import {catchError, Observable, switchMap, throwError} from "rxjs";
import {inject} from "@angular/core";
import {AuthService} from "../services/auth-service/auth.service";

let isRefreshing = false

export function authInterceptor(req: HttpRequest<unknown>, next: HttpHandlerFn): Observable<HttpEvent<unknown>> {
  const authService: AuthService = inject(AuthService);
  const token: string | null = authService.getToken;

  if (!authService.isAuth || token == null) {
    return next(req);
  }

  if (isRefreshing) {
    return refreshAndProceed(authService, req, next);
  }

  return next(addToken(req, token))
    .pipe(
      catchError(error => {
        if (error.status === 401) {
          return refreshAndProceed(authService, req, next);
        }
        return throwError(error);
      })
    )
}


const refreshAndProceed = (
  authService: AuthService,
  req: HttpRequest<unknown>,
  next: HttpHandlerFn
) => {
  if (!isRefreshing) {
    isRefreshing = true
    return authService.refreshToken()
      .pipe(
        switchMap((response) => {
          isRefreshing = false
          return next(addToken(req, response.data));
        })
      )
  }
  return next(addToken(req, authService.token!));
}

const addToken = (req: HttpRequest<unknown>, token: string) => {
  return req.clone({
    setHeaders: {
      Authorization: `Bearer ${token}`
    }
  })
}
