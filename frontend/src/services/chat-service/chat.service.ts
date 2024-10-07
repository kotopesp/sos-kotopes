import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, forkJoin, map, Observable, of, switchMap, throwError } from 'rxjs';
import { Chat, ResponseUser } from '../../model/chat.interface'
import { environment } from '../../environments/environment'
import { User } from '../../model/user.interface';
import { Message } from '../../model/message.interface';
import { AuthService } from '../auth-service/auth.service';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = environment.apiUrl;
  private chatsSubject: BehaviorSubject<Chat[]> = new BehaviorSubject<Chat[]>([]);
  chats$: Observable<Chat[]> = this.chatsSubject.asObservable();

  constructor(private http: HttpClient, private authService: AuthService) { }

  getById(id: number): Observable<Chat> {
    const url = `${this.apiUrl}chats/${id}`;
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    const chat = this.http.get<Chat>(url, {headers}).pipe(
      map((chat) => ({
        ...chat,
        unread_count: 0,
      } as Chat)
    ));
    this.readMessages(id);
    return chat;
  }

  selectedUserIds: number[] = [];

  createChat(selectedUserIds: number[], userId: number): Observable<Chat> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    return this.http.post<{data: Chat}>(`${this.apiUrl}chats`, { userIds: selectedUserIds}, {headers}).pipe(
      catchError((error: HttpErrorResponse) => {
        // 409 (Conflict), перенаправляем на существующий чат
        if (error.status === 409 && error.error && error.error.data && error.error.data.id) {
          const resp = this.getChatById(error.error.data.id).pipe(
            map(chat => (
              {data: chat}
            ))
          );
          return resp;
        } else {
          return throwError(() => error); // Если ошибка другая, пробрасываем дальше
        }
      }),
      map(chat => ({
        ...chat.data,
        title: this.getTitle(chat.data.users, userId),
        unread_count: 0,
      })
    ));
  }

  updateChat(message: Message, currentChat: Chat) {
    const chats = this.chatsSubject.value;
    const chatIndex = chats.findIndex(c => c.id === message.chat_id);

    if (chatIndex >= 0) {
      chats[chatIndex].last_message = {
          message_content: message.message_content,
          created_at: message.created_at,
          user_id: message.user_id,
          is_read: message.user_id === this.authService.getIdFromToken,
          sender_name: message.user_id === this.authService.getIdFromToken ? "Вы" : message.sender_name,
        };
      if (chats[chatIndex].id === currentChat.id) {
        this.readMessages(currentChat.id);
      }
    }

    this.chatsSubject.next([...chats]);
  }

  getFavUsers(): Observable<{ id: number, username: string }[]> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    return this.http.get<ResponseUser[]>(`${this.apiUrl}users/favourites`, {headers})
      .pipe(
        map(users => users.map(user => ({
          id: user.id,
          username: user.username,
        })))
      );
  }

  getTitle(users: User[], currentUser: number) : string {
    const sortusers = users.sort().filter(u => u.id != currentUser)
    return sortusers.length != 0 ? sortusers.map(user => user.username).join(', ') : "Избранное";
  }

  getChatById(chatId: number): Observable<Chat> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    const resp = this.http.get<{ data: Chat, }>(`${this.apiUrl}chats/${chatId}`, {headers})
    .pipe(
      map(responce => ({
        ...responce.data,
        unread_count: 0,
      })
    ));
    return resp;
  }

  readMessages(chatId: number) {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    const resp = this.http.patch(`${this.apiUrl}chats/${chatId}/unread`, null, {headers}).pipe(
      map(resp => {
        return resp;
      })
    );
    resp.subscribe(
      resp => resp
    );
  }

  getUnreadCount(chatId: number): Observable<number> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    return this.http.get<{ data: number }>(`${this.apiUrl}chats/${chatId}/unread`, { headers}).pipe(
      map(response => {
          return response.data;
        },
      ));
  }

  getMessagesByChatId(chatId: number): Observable<Message[]> {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    return this.http.get<{ data: { message: Message[] } }>(`${this.apiUrl}chats/${chatId}/messages`, { headers}).pipe(
      map(response => {
        const messages = response.data.message || [];
        return messages.map(message => {
          message.is_user_message = message.user_id === this.authService.getIdFromToken;
          message.time = (new Date(message.created_at)).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
          return message;
        });
      }),
      catchError(() => of([])) 
    );
  }

  getAllChats(userId: number) {
    const headers = new HttpHeaders({
      'Authorization': `Bearer ${this.authService.getToken}`
    });
    this.http.get<{ data: { chats: Chat[], total: number } }>(`${this.apiUrl}chats`, { headers} )
      .pipe(
        switchMap(response => {
          const chats = response.data.chats;
          const chatObservables = chats.map(chat => 
            this.getUnreadCount(chat.id).pipe(
              map(unreadCount => ({
                ...chat,
                title: this.getTitle(chat.users, userId),
                last_message: {
                  ...chat.last_message, 
                  sender_name: chat.last_message ? (chat.last_message.user_id == userId ? "Вы" : chat.last_message.sender_name) : "",
                },
                unread_count: unreadCount
              } as Chat)),
              catchError(() => of({
                ...chat,
                title: this.getTitle(chat.users, userId),
                unread_count: 0
              }))
            )
          );
          
          return forkJoin(chatObservables);
        }),
        catchError(() => of([]))
      )
      .subscribe(chats => {
        chats.sort((a: Chat, b: Chat) => 
          {
            const aLastMessageDate = a.last_message ? new Date(a.last_message.created_at).getTime() : new Date(a.created_at).getTime();
            const bLastMessageDate = b.last_message ? new Date(b.last_message.created_at).getTime() : new Date(b.created_at).getTime();
        
            return bLastMessageDate - aLastMessageDate;
          });
        this.chatsSubject.next(chats);
      });
    }
}