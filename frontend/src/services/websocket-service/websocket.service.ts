import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { map, Observable, Subject } from 'rxjs';
import { environment } from '../../environments/environment'
import { HttpClient } from '@angular/common/http';
import { Message } from '../../model/message.interface';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private socket!: WebSocket;
  private messagesSubject: Subject<string> = new Subject<string>();

  constructor(private http: HttpClient) {}
  private apiUrl = environment.apiUrl;
  private chatId = -1;
  connect(chatID: number): void {
    this.socket = new WebSocket(`ws:${this.apiUrl.split("//")[1]}ws/${chatID}`);
    this.chatId = chatID;
    console.log("lalalala", this.socket);
    this.socket.onopen = () => {
      console.log('Connected to WebSocket');
    };

    this.socket.onmessage = (event) => {
      console.log('Message from server: ', event);
      this.messagesSubject.next(event.data);
    };

    this.socket.onclose = () => {
      console.log('WebSocket connection closed');
    };
  }

  sendMessage(message: string): void {
    console.log(this.socket);
    if (this.socket.readyState === WebSocket.OPEN) {
      var resp = this.http.post<{data: any}>(`${this.apiUrl}chats/${this.chatId}/messages`, JSON.parse(message)[0]).pipe(
              map(response => <Message>{
                ID: response.data.ID,
                UserID: response.data.UserID,
                ChatID: this.chatId,
                Content: response.data.Content,
              }),
            );
            resp.subscribe(
              (response => response),
            );
      this.socket.send(message);
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
      console.log('WebSocket connection closed');
    }
  }
}

// export class WebsocketService {
//   private socket$: WebSocketSubject<string>;
  
//   private apiUrl = environment.apiUrl;

// //  constructor(id: number) {
//     constructor(private http: HttpClient) {
//       this.socket$ = webSocket('ws://localhost:8080/chats/ws');
//       //this.socket$ = webSocket(`ws://${this.apiUrl}chats/ws`);
//   }

//   // Method to send messages to the server
//   public sendMessage(msg: string, chatId: number): void {
//     var resp = this.http.post<{data: any}>(`${this.apiUrl}chats/${chatId}/messages`, JSON.parse(msg)[0]).pipe(
//       map(response => <Message>{
//         ID: response.data.ID,
//         UserID: response.data.UserID,
//         ChatID: chatId,
//         Content: response.data.Content,
//       }),
//     );
//     resp.subscribe(
//       (response => response),
//     );
//     this.socket$.next(msg);
//   }

//   // Method to receive messages from the server
//   public getMessages(): Observable<string> {
//     return this.socket$.asObservable();
//   }
// }
