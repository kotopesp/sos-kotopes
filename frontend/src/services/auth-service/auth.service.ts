import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Observable, tap} from "rxjs";
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
    return this.cookieService.get("token");
  }

  login(payload: {
    password: string,
    username: string,
  }): Observable<LoginResponse> {
    return this.http.post<LoginResponse>(
      `${this.baseApiUrl}auth/login`,
      payload,
      {withCredentials: true}
    ).pipe(
      tap((res: LoginResponse) => {
          this.saveTokens(res);
        }
      )
    );
  }

  logout() {
    this.cookieService.deleteAll()
    this.token = null;
    this.router.navigate([''])
  }

  setToken(token: string) {
    this.token = token;
  }

  refreshToken() {
    this.http.post<LoginResponse>(
      `${this.baseApiUrl}auth/token/refresh`,
      {},
      {withCredentials: true}
    ).subscribe((res: LoginResponse) => {
      this.saveTokens(res);
    })
  }

  saveTokens(res: LoginResponse) {
    this.setToken(res.data);
    this.cookieService.set('token', res.data)
  }
}
