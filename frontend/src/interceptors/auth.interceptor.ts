import {HttpEvent, HttpHandlerFn, HttpRequest} from "@angular/common/http";
import {Observable} from "rxjs";
import {inject} from "@angular/core";
import {AuthService} from "../services/auth-service/auth.service";

export function authInterceptor(req: HttpRequest<unknown>, next: HttpHandlerFn): Observable<HttpEvent<unknown>> {
  const authService = inject(AuthService);
  const token: string | null = authService.getToken;

  if (authService.isAuth && token != null) {
    const newReq = req.clone({
      headers: req.headers.append('Authorization', "Bearer " + token),
    });
    return next(newReq);
  }

  return next(req);
}
