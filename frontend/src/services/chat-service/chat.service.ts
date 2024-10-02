import { HttpClient, HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { BehaviorSubject, catchError, concatMap, firstValueFrom, forkJoin, map, Observable, of, OperatorFunction, switchMap, tap, throwError } from 'rxjs';
import { Chat, ResponseUser } from '../../model/chat.interface'
import { environment } from '../../environments/environment'
import { User } from '../../model/user.interface';
import { Message } from '../../model/message.interface';

@Injectable({
  providedIn: 'root'
})
export class ChatService {
  private apiUrl = environment.apiUrl;
  private chatsSubject: BehaviorSubject<Chat[]> = new BehaviorSubject<Chat[]>([]);
  chats$: Observable<Chat[]> = this.chatsSubject.asObservable();

  constructor(private http: HttpClient) { }

  getById(id: number): Observable<Chat> {
    const url = `${this.apiUrl}chats/${id}`;
    var chat = this.http.get<Chat>(url).pipe(
      map((chat) => <Chat>{
        ...chat,
        unread_count: 0,
      })
    );
    this.readMessages(id);
    return chat;
  }

  selectedUserIds: number[] = [];

  createChat(selectedUserIds: number[], userId: number): Observable<Chat> {
    return this.http.post<{data: Chat}>(`${this.apiUrl}chats`, { userIds: selectedUserIds }).pipe(
      catchError((error: HttpErrorResponse) => {
        // 409 (Conflict), перенаправляем на существующий чат
        if (error.status === 409 && error.error && error.error.data && error.error.data.id) {
          var resp = this.getChatById(error.error.data.id).pipe(
            map(chat => (
              {data: chat}
            ))
          );
          return resp;
        } else {
          return throwError(() => error); // Если ошибка другая, пробрасываем дальше
        }
      }),
      map(chat => <Chat>{
        ...chat.data,
        title: this.getTitle(chat.data.users, userId),
        unread_count: 0,
      })
    );
  }

  updateChat(message: Message, currentUserId: number) {
    const chats = this.chatsSubject.value;
    const chatIndex = chats.findIndex(c => c.id === message.chat_id);

    if (chatIndex >= 0) {
      chats[chatIndex].last_message = {
          message_content: message.message_content,
          created_at: message.created_at,
          user_id: message.user_id,
          is_read: message.is_user_message,
          sender_name: message.is_user_message ? "Вы" : message.sender_name,
        };
      chats[chatIndex].unread_count = (message.user_id === currentUserId) ? 0 : chats[chatIndex].unread_count + 1;
    }

    this.chatsSubject.next([...chats]);
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

  getTitle(users: User[], currentUser: number) : string {
    var sortusers = users.sort().filter(u => u.id != currentUser) // todo подправить отображение заголовка в чате
    return sortusers.map(user => user.username).join(', ');
  }

  getChatById(chatId: number): Observable<Chat> {
    var resp = this.http.get<{ data: any, }>(`${this.apiUrl}chats/${chatId}`)
    .pipe(
      map(responce => <Chat> {
        id: responce.data.id,
        chat_type: responce.data.chat_type,
        users: responce.data.users,
        unread_count: 0,
      })
    );
    return resp;
  }

  readMessages(chatId: number) {
    var resp = this.http.patch(`${this.apiUrl}chats/${chatId}/unread/`, null).pipe(
      map(resp => {
        const chats = this.chatsSubject.value;
        const chatIndex = chats.findIndex(chat => chat.id === chatId);

        if (chatIndex >= 0) {
          chats[chatIndex].unread_count = 0;
          this.chatsSubject.next([...chats]);
        }

        return resp;
      })
    );
    resp.subscribe(
      resp => resp
    );
  }

  getUnreadCount(chatId: number): Observable<number> {
    return this.http.get<{ data: number }>(`${this.apiUrl}chats/${chatId}/unread`).pipe(
      map(response => {
          return response.data;
        },
      ));
  }

  getMessagesByChatId(chatId: number, userId: number): Observable<Message[]> {
    return this.http.get<{ data: { message: Message[] } }>(`${this.apiUrl}chats/${chatId}/messages`).pipe(
      map(response => {
        const messages = response.data.message || [];
        return messages.map(message => {
          message.is_user_message = message.user_id === userId;
          message.time = (new Date(message.created_at)).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
          return message;
        });
      }),
      catchError(() => of([])) 
    );
  }

  getAllChats(userId: number) {
    this.http.get<{ data: { chats: Chat[], total: number } }>(`${this.apiUrl}chats`)
      .pipe(
        switchMap(response => {
          const chats = response.data.chats;
  
          const chatObservables = chats.map(chat => 
            this.getUnreadCount(chat.id).pipe(
              map(unreadCount => (<Chat>{
                ...chat,
                title: this.getTitle(chat.users, userId),
                last_message: {
                  ...chat.last_message, 
                  sender_name: chat.last_message ? (chat.last_message.user_id == userId ? "Вы" : chat.last_message.sender_name) : "",
                },
                unread_count: unreadCount
              })),
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
        this.chatsSubject.next(chats);
      });
    }
}