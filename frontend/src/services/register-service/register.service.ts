import {inject, Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {FormControl, ɵFormGroupValue, ɵTypedOrUntyped} from "@angular/forms";

@Injectable({
  providedIn: 'root'
})
export class RegisterService {


  http = inject(HttpClient);
  baseApiUrl: string = environment.apiUrl;

  registration(payload: ɵTypedOrUntyped<{ password: FormControl<null>; username: FormControl<null> }, ɵFormGroupValue<{
    password: FormControl<null>;
    username: FormControl<null>
  }>, any>) {
    return this.http.post(
      `${this.baseApiUrl}auth/signup`,
      payload
    )
  }
}
