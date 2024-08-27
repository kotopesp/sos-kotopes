import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {FormControl, ɵFormGroupValue, ɵTypedOrUntyped} from "@angular/forms";
import {Observable} from "rxjs";

interface RegistrationResponse {
  token: string;
  user: {
    id: number;
    name: string;
    email: string;
  };
}

@Injectable({
  providedIn: 'root'
})
export class RegisterService {


  http = inject(HttpClient);
  baseApiUrl: string = environment.apiUrl;

  registration(payload: ɵTypedOrUntyped<{ password: FormControl<null>; username: FormControl<null> }, ɵFormGroupValue<{
    password: FormControl<null>;
    username: FormControl<null>
  }>, string>): Observable<RegistrationResponse> {
    return this.http.post(
      `${this.baseApiUrl}auth/signup`,
      payload
    );
  }
}
