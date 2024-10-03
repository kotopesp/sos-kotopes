import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Observable} from "rxjs";
import {User} from "../../model/user.interface";

@Injectable({
  providedIn: 'root'
})
export class RegisterService {


  http = inject(HttpClient);
  baseApiUrl: string = environment.apiUrl;

  registration(payload: {
    email: string,
    name: string,
    lastname: string,
    username: string,
    password: string
  }): Observable<User> {
    return this.http.post(
      `${this.baseApiUrl}auth/signup`,
      payload
    )
  }
}
