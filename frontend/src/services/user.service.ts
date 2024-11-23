import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { User } from '../model/user.interface'
import { environment } from '../environments/environment'

@Injectable({
  providedIn: 'root'
})
export class UserService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

    getById(id: string | null): Observable<User> {
    return this.http.get<User>(`${this.apiUrl}/users/${id}`)
  }

  update(user: User): Observable<User> {
    return this.http.patch(`${this.apiUrl}/users`, user)
  }
}
