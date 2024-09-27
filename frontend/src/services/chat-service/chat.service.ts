import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, concatMap, firstValueFrom, map, Observable, of, OperatorFunction, tap, throwError } from 'rxjs';
import { Chatmember } from '../../model/chatmember.interface'
import { Chatinfo } from '../../model/chatinfo.interface'
import { Chat, ResponseUser } from '../../model/chat.interface'
import { environment } from '../../environments/environment'
import { User } from '../../model/user.interface';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  getById(id: number): Observable<Chat> {
    const url = `${this.apiUrl}chats/${id}`;
    var chat = this.http.get<Chat>(url).pipe(
      map((chat: Chat) => chat)
    );
    return chat;
  }

  selectedUserIds: number[] = []; // ID выбранных пользователей

  createChat(selectedUserIds: number[]): Observable<Chatinfo> {
    return this.http.post<{data: Chat}>(`${this.apiUrl}chats`, { userIds: selectedUserIds }).pipe(
      catchError((error: HttpErrorResponse) => {
        // 409 (Conflict), перенаправляем на существующий чат
        if (error.status === 409 && error.error && error.error.data && error.error.data.ID) {
          return this.getChatById(error.error.data.ID).pipe(
            map(chat => ({data: chat}))
          );
        } else {
          return throwError(() => error); // Если ошибка другая, пробрасываем дальше
        }
      }),
      map(chat => <Chatinfo>{
        Id: chat.data.ID,
        Chattype: chat.data.ChatType,
        Title: this.getTitle(chat.data.Users), // chat уже обработан как тип Chat
      })
    );
  }

  getFavUsers(): Observable<{ id: number, username: string }[]> {
    return this.http.get<ResponseUser[]>(`${this.apiUrl}users/favourites`)
      .pipe(
        map(users => users.map(user => ({
          id: user.id,
          username: user.username,
        })))
      );
  }

  getTitle(users: User[]) : string {
    // if length(chat.users) > 1{
    //   return chat.title
    // }
    var title = ""
    for (var user of users.sort()) {
      title += user.Username + ", "
    }
    return title.slice(0, title.length - 2);
  }

  getChatById(chatId: number): Observable<Chat> {
    return this.http.get<{ data: any, }>(`${this.apiUrl}chats/${chatId}`)
    .pipe(
      map(responce => <Chat> {
        ID: responce.data.ID,
        ChatType: responce.data.ChatType,
        Users: responce.data.Users,
      })
    );
  }

  getAllChats(): Observable<Chatinfo[]> {
    return this.http.get<{ data: { chats: Chat[], total: number } }>(`${this.apiUrl}chats`)
    .pipe(
      map(responce => responce.data.chats.map(chat => (<Chatinfo>{
        Id: chat.ID,
        Chattype: chat.ChatType,
        Title: this.getTitle(chat.Users),
      })))
    );
  }
}
