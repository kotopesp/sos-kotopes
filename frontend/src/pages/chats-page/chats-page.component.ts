import {AfterViewChecked, Component, ElementRef, signal, ViewChild, OnInit} from '@angular/core';
import {AppChatTypeButtonComponent} from "../../entities/chat-type-button/app-chat-type-button.component";
import {Button} from "../../model/button";
import {NgClass, NgForOf, NgIf, NgStyle, NgSwitch, NgSwitchCase} from "@angular/common";
import {ChatComponent} from "./ui/chat/chat.component";
import {MessageComponent} from "./ui/message/message.component";
import {PostAnswerComponent} from "./ui/post-answer/post-answer.component";
import {AddToChatComponent} from "./ui/add-to-chat/add-to-chat.component";
import {ToggleActiveDirective} from "./toggle-active.directive";
import { WebsocketService } from '../../services/websocket-service/websocket.service';
import { FormControl, FormGroup, Validators, ReactiveFormsModule  } from '@angular/forms';

@Component({
  selector: 'app-chats-page',
  standalone: true,
  imports: [
    AppChatTypeButtonComponent,
    NgForOf,
    NgIf,
    ChatComponent,
    MessageComponent,
    PostAnswerComponent,
    AddToChatComponent,
    NgSwitch,
    NgSwitchCase,
    NgStyle,
    ToggleActiveDirective,
    NgClass,
    ReactiveFormsModule,
  ],
  templateUrl: './chats-page.component.html',
  styleUrl: './chats-page.component.scss'
})
export class ChatsPageComponent implements AfterViewChecked, OnInit {
  currentChat = false;
  createChat = false;

  countInArray = signal<number>(0);
  @ViewChild('scrollableContainer', { static: false }) private scrollableContainer?: ElementRef;

  buttons: Button[] = [
    {label: "Все чаты", color: "#352B1A", iconName: "allChats.svg"},
    {label: "Отклики", color: "#352303", counter: 5, iconName: "respond.svg"},
    {label: "Чаты с передержкой", color: "#2B1800", counter: 4, iconName: "keepers.svg"},
    {label: "Чаты с отловщиками", color: "#221630", iconName: "seekers.svg"},
    {label: "Чаты с ветеринарами", color: "#182C2A", iconName: "vets.svg"},
    {label: "Другие чаты", color: "#21190B", iconName: "other.svg"},
  ]

  sendMsgForm : FormGroup = new FormGroup({
    "msgText": new FormControl("", [
                Validators.required,
    ]) 
});

  ngAfterViewChecked(): void {
    this.scrollToBottom();
  }

  private scrollToBottom(): void {
    if (this.scrollableContainer && this.scrollableContainer.nativeElement) {
      const container = this.scrollableContainer.nativeElement as HTMLElement;
      container.scrollTop = container.scrollHeight;
    }
  }

  public messages: string[] = [];
  public messageText: string = '';

  constructor(private websocketService: WebsocketService) {}

  ngOnInit(): void {
    // Subscribe to incoming messages from WebSocket
    this.websocketService.getMessages().subscribe((msg: string) => {
      this.messages.push(msg);
    });
  }

  // onMessageInput(event: Event): void {
  //   const target = event.target as HTMLTextAreaElement;
  //   this.messageText = target.value;  // Обновляем текст сообщения
  // }

  // Method to send a message to the WebSocket server
  onSubmit() {
    this.messageText = this.sendMsgForm.controls['msgText'].value;
    this.sendMsgForm.reset();
    if (this.messageText.trim()) {
      this.websocketService.sendMessage(this.messageText);
      this.messageText = '';
    }
  }

  isUserMessage(message: { content: string, userId: string }): boolean {
    return true; //message.userId === this.currentUserId;  // сравниваем userId сообщения с текущим пользователем
  }

}