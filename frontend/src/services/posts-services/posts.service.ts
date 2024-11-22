import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import {Router} from "@angular/router";
import {Post, PostResponse} from "../../model/post.interface";
import {map, Observable} from "rxjs";
import {User} from "../../model/user.interface";


@Injectable({
  providedIn: 'root'
})
export class PostsService {
  limit: string | null = null;
  offset: string | null = null;

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


  getPosts(): Observable<PostResponse> {
    const params = { limit: this.limit || '10', offset: this.offset || '0' };

    return this.http.get<PostResponse>(`${environment.apiUrl}posts`, { params});
  }
}
