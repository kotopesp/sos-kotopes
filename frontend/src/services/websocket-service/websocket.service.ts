import { Injectable } from '@angular/core';
import { WebSocketSubject, webSocket } from 'rxjs/webSocket';
import { Observable } from 'rxjs';
import { environment } from '../../environments/environment'

@Injectable({
  providedIn: 'root'
})
export class WebsocketService {
  private socket$: WebSocketSubject<string>;
  
  private apiUrl = environment.apiUrl;

//  constructor(id: number) {
    // Open WebSocket connection to Go WebSocket server
    constructor() {
      this.socket$ = webSocket('ws://localhost:8080/chats/ws');
      //this.socket$ = webSocket(`ws://${this.apiUrl}chats/ws`);
      console.log(`${this.apiUrl}chats`)
  }

  // Method to send messages to the server
  public sendMessage(msg: string): void {
    this.socket$.next(msg);
    console.log("send message", msg); // TODO: delete
  }

  // Method to receive messages from the server
  public getMessages(): Observable<string> {
    return this.socket$.asObservable();
  }
}
