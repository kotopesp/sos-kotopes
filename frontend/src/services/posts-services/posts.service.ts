import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Router} from "@angular/router";
import {FormGroup} from "@angular/forms";


@Injectable({
  providedIn: 'root'
})
export class PostsService {

  constructor(private http: HttpClient, private router: Router) {
  }

  createPost(payload: FormGroup) {
    return this.http.post<any>(`${environment.apiUrl}/posts`, payload)
  }
}
