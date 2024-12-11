import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { User } from '../../model/user.interface'
import { environment } from '../../environments/environment'

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  getById(id: number): Observable<User> {
    return this.http.get(`${this.apiUrl}/users/${id}`).pipe(
      map((user:User) => user)
    )
  }

  update(user: User): Observable<User> {
    return this.http.put(`${this.apiUrl}/users/${user.id}`, user)
  }
}
