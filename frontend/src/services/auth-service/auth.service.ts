import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Observable, tap} from "rxjs";
import {CookieService} from "ngx-cookie-service";
import {Router} from "@angular/router";


export interface LoginResponse {
  status: string,
  data: {
    access_token: string
  }
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

  get userID() {
    const decodedData = this.decodeToken(this.cookieService.get('token'));
     // Извлекаем id пользователя
    return decodedData.id
  }

  get Token() {
    return this.cookieService.get('token')
  }

  get isAuth() {
    if (!this.token) {
      this.token = this.cookieService.get('token')
    }
    return !!this.token;
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
          this.saveTokens(res)
        }
      )
    );
  }

  logout() {
    this.cookieService.deleteAll();
    this.token = null;
    this.router.navigate([''])
  }

  setToken(token: string) {
    this.token = token;
  }

  saveTokens(res: LoginResponse) {
    this.setToken(res.data.access_token);
    this.cookieService.set('token', res.data.access_token)
  }

  decodeToken(token: string): any {
    // JWT состоит из трех частей, разделенных точками. Получаем полезную нагрузку (payload).
    const payload = token.split('.')[1];
    // Декодируем строку Base64
    const decodedPayload = atob(payload);
    // Преобразуем декодированную строку в объект JSON
    return JSON.parse(decodedPayload);
  }
}
