import {HttpHandlerFn, HttpInterceptorFn, HttpRequest} from "@angular/common/http";
import {AuthService, LoginResponse} from "../auth-service/auth.service";
import {inject} from "@angular/core";
import {catchError, switchMap, throwError} from "rxjs";

let isRefreshing = false;

export const authTokenInterceptor: HttpInterceptorFn = (req, next) => {
  const authService = inject(AuthService);
  const token = authService.token;

  if (!token) return next(req)

  if (isRefreshing) {
    return refreshOldToken(authService, req, next);
  }

  return next(addToken(req, token))
    .pipe(
      catchError(error => {
        if (error.status === 403) {
          return refreshOldToken(authService, req, next);
        }
        return throwError(error)
      })
    )
}

export const refreshOldToken = (
  authService: AuthService,
  req: HttpRequest<any>,
  next: HttpHandlerFn
) => {
  if (!isRefreshing) {
    isRefreshing = true;
    return authService.refreshAuthToken()
      .pipe(
        switchMap((res: LoginResponse) => {
          isRefreshing = false;
          return next(addToken(req, res.data.access_token))
        })
      )
  }
  return next(addToken(req, authService.token!))
}

export const addToken = (req: HttpRequest<any>, token: string) => {
  return req.clone({
    setHeaders: {
      Authorization: `Bearer ${token}`
    }
  })
}
