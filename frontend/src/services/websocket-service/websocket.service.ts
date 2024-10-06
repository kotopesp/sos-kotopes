import { Injectable } from '@angular/core';
import { map, Observable, Subject } from 'rxjs';
import { environment } from '../../environments/environment'
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Message } from '../../model/message.interface';
import { AuthService } from '../auth-service/auth.service';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private socket!: WebSocket;
  private messagesSubject: Subject<string> = new Subject<string>();

  constructor(private http: HttpClient, private authService: AuthService) {}
  private apiUrl = environment.apiUrl;
  private chatId = -1;
  connect(chatID: number): void {
    this.socket = new WebSocket(`ws:${this.apiUrl.split("//")[1]}ws/${chatID}`);
    this.chatId = chatID;
    this.socket.onopen = () => {
    };

    this.socket.onmessage = (event) => {
      this.messagesSubject.next(event.data);
    };

    this.socket.onclose = () => {
    };
  }

  sendMessage(message: string): void {
    if (this.socket.readyState === WebSocket.OPEN) {
      const headers = new HttpHeaders({
        'Authorization': `Bearer ${this.authService.getToken}`
      });
      var resp = this.http.post<{data: any}>(`${this.apiUrl}chats/${this.chatId}/messages`, JSON.parse(message)[0], {headers}).pipe(
              map(response => <Message>{
                user_id: response.data.user_id,
                chat_id: this.chatId,
                message_content: response.data.message_content,
                sender_name: response.data.sender_name,
                created_at: response.data.created_at,
              }),
            );
            resp.subscribe(
              (response => this.socket.send(JSON.stringify([response]))),
            );
    } else {
      console.error('WebSocket is not open.');
    }
  }

  getMessages(): Observable<string> {
    return this.messagesSubject.asObservable();
  }

  closeConnection(): void {
    if (this.socket) {
      this.socket.close();
    }
  }
}
