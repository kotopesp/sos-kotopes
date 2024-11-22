import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Router} from "@angular/router";


@Injectable({
  providedIn: 'root'
})
export class PostsService {

  constructor(private http: HttpClient, private router: Router) {
  }


  createPost(payload: FormData) {
    return this.http.post<any>(`${environment.apiUrl}posts`, payload).subscribe(
      {
        next: () => {
          console.log('success')
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }


  getPosts() {
    return this.http.get(`${environment.apiUrl}posts`);
  }
}
