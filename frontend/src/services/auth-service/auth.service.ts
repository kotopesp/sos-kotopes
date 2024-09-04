import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {FormControl, ɵFormGroupValue, ɵTypedOrUntyped} from "@angular/forms";
import {Observable, tap} from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  private token = null;

  http = inject(HttpClient);
  baseApiUrl: string = environment.apiUrl;

  login(payload: {email_or_login: string, password: string}):
    Observable<{token: string}>{
    return this.http.post<{token: string}>(
      `${this.baseApiUrl}auth/login`,
      payload
    ).pipe(
      tap(
        ({token}) => {
          localStorage.setItem('auth-token', token)
          this.setToken(token);
        }
      )
    );
  }

  setToken(token: string) {
    this.token = token;
  }
}
