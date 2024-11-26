import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {environment} from "../../environments/environment";
import { PostResponse} from "../../model/post.interface";
import { Observable} from "rxjs";


@Injectable({
  providedIn: 'root'
})
export class PostsService {
  limit: string | null = null;
  offset: string | null = null;

  constructor(private http: HttpClient) {
  }


  createPost(payload: FormData) {
    return this.http.post<any>(`${environment.apiUrl}posts`, payload).subscribe(
      {
        next: () => {
          console.log('success') // TO DO: add routing to success page (overlay)
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

  getPostsUser(user_id: string | null): Observable<PostResponse> {
    const params = { limit: this.limit || '10', offset: this.offset || '0' };
    return this.http.get<PostResponse>(`${environment.apiUrl}users/${user_id}/posts`, { params})
  }

  getPostsFavoritesUser(): Observable<PostResponse> {
    const params = { limit: this.limit || '10', offset: this.offset || '0' };
    return this.http.get<PostResponse>(`${environment.apiUrl}posts/favourites`, { params})
  }

  addPostToFavorites(post_id: number | null) {
    return this.http.post<PostResponse>(`${environment.apiUrl}posts/${post_id}/favourites`, post_id).subscribe(
      {
        next: () => {
          console.log('success') // TO DO: add routing to success page (overlay)
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }

  deletePostFromFavorites(post_id: number | null) {
    return this.http.delete<PostResponse>(`${environment.apiUrl}posts/${post_id}/favourites`).subscribe(
      {
        next: () => {
          console.log('success') // TO DO: add routing to success page (overlay)
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }
}