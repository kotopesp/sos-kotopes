import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {catchError, Observable, tap, throwError} from "rxjs";
import {CookieService} from "ngx-cookie-service";
import {Router} from "@angular/router";


export interface LoginResponse {
  status: string,
  data: string
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {

  token: string | null;
  baseApiUrl: string;

  constructor(private http: HttpClient, private router: Router, private cookieService: CookieService) {
    this.token = null;
    this.baseApiUrl = environment.apiUrl;
  }

  get isAuth() {
    if (!this.token) {
      this.token = this.cookieService.get('token')
    }
    return !!this.token;
  }

  get getToken() : string | null {
    return localStorage.getItem('authToken');
  }

  get getIdFromToken() : number {
    if (!this.token) {
      this.token = this.cookieService.get('token')
    }
    const payload = this.token.split('.')[1];
    const decoded = atob(payload);
    const decObject = JSON.parse(decoded);
    const res = decObject.id;
    return res;
  }

  login(payload: {
    password: string,
    username: string,
  }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(
      `${this.baseApiUrl}auth/login`,
      payload
    ).pipe(
      tap((res: LoginResponse) => {
          this.saveTokens(res);
          localStorage.setItem('authToken', res.data);
        }
      )
    );
  }

  refreshAuthToken() {
    return this.http.post(
      `${this.baseApiUrl}auth/v1/auth/token/refresh`,
      ''
    ).pipe(catchError(
        error => {
          this.logout()

          return throwError(error);
        }
      )
    )
  }

  logout() {
    this.cookieService.deleteAll()
    this.token = null;
    this.router.navigate([''])
  }

  setToken(token: string) {
    this.token = token;
  }

  saveTokens(res: LoginResponse) {
    this.setToken(res.data);
    this.cookieService.set('token', res.data)
  }
}
