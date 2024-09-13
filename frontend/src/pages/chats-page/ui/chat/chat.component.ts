import {Component, OnInit} from '@angular/core';
import {NgClass} from "@angular/common";
import {ToggleActiveDirective} from "../../toggle-active.directive";

import { WebsocketService } from '../../../../services/websocket-service/websocket.service';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [
    NgClass,
    ToggleActiveDirective
  ],
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.scss'
})
export class ChatComponent implements OnInit {
  public messages: string[] = [];
  public message: string = '';

  constructor(private websocketService: WebsocketService) {}

  ngOnInit(): void {
    // Subscribe to incoming messages from WebSocket
    this.websocketService.getMessages().subscribe((msg: string) => {
      this.messages.push(msg);
    });
  }

  // Method to send a message to the WebSocket server
  sendMessage(): void {
    if (this.message.trim()) {
      this.websocketService.sendMessage(this.message);
      this.message = '';
    }
  }
}
