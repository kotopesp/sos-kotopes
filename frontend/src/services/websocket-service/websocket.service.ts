import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { map, Observable } from 'rxjs';
import { environment } from '../../environments/environment'
import { HttpClient } from '@angular/common/http';
import { Message } from '../../model/message.interface';

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private socket$: WebSocketSubject<string>;
  
  private apiUrl = environment.apiUrl;

//  constructor(id: number) {
    constructor(private http: HttpClient) {
      this.socket$ = webSocket('ws://localhost:8080/chats/ws');
      //this.socket$ = webSocket(`ws://${this.apiUrl}chats/ws`);
  }

  // Method to send messages to the server
  public sendMessage(msg: string, chatId: number): void {
    var resp = this.http.post<{data: any}>(`${this.apiUrl}chats/${chatId}/messages`, JSON.parse(msg)[0]).pipe(
      map(response => <Message>{
        ID: response.data.ID,
        UserID: response.data.UserID,
        ChatID: chatId,
        Content: response.data.Content,
      }),
    );
    resp.subscribe(
      (response => response),
    );
    this.socket$.next(msg);
  }

  // Method to receive messages from the server
  public getMessages(): Observable<string> {
    return this.socket$.asObservable();
  }
}
