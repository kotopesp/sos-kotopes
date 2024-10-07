import {AfterViewChecked, Component, ElementRef, signal, ViewChild, OnInit, OnDestroy} from '@angular/core';
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
import { map, switchMap } from 'rxjs';
import { ActivatedRoute, Params, Router } from '@angular/router';
import {ChatService} from '../../services/chat-service/chat.service';
import { HttpClient } from '@angular/common/http';
import { Chat } from '../../model/chat.interface';
import { Message } from '../../model/message.interface';
import { AuthService } from '../../services/auth-service/auth.service';

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
export class ChatsPageComponent implements AfterViewChecked, OnInit, OnDestroy {
  currentChat: Chat = { id: -1, title: '', chat_type: 'other', unread_count: 0 , users: [], created_at: new Date};
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
    "msgText": new FormControl("", [Validators.required,]) 
  });

  private previousMessageCount = 0;

  ngAfterViewChecked(): void {
    if (this.messages.length !== this.previousMessageCount) {
      this.scrollToBottom();
      this.previousMessageCount = this.messages.length;
    }
  }

  private scrollToBottom(): void {
    if (this.scrollableContainer && this.scrollableContainer.nativeElement) {
      const container = this.scrollableContainer.nativeElement as HTMLElement;
      container.scrollTop = container.scrollHeight;
    }
  }

  public messages: Message[] = [];
  public messageText = '';
  public favusers: { id: number, username: string}[] = [];
  public chatList: Chat[] = [];
  public userId = -1;
  private refreshInterval!: NodeJS.Timeout;
  
  constructor(private router: Router, private authService: AuthService, private activatedRoute: ActivatedRoute, private chatService: ChatService, private websocketService: WebsocketService, private http: HttpClient) {}
  
  ngOnInit(): void {
    this.activatedRoute.params.pipe(
      map((params: Params) => parseInt(params['id'], 10)),
      switchMap((chatId: number) => this.chatService.getById(chatId))
    );
    
    this.websocketService.getMessages().subscribe((msg: string) => {
      const message = JSON.parse(msg)[0] as Message;
      this.chatService.updateChat(message, this.currentChat);
      this.messages.push(message);
    });
    
    this.userId = this.authService.getIdFromToken;
    this.chatService.getFavUsers().subscribe(
      (users) => {
        this.favusers = users;
      }
    );
    
    this.chatService.getAllChats(this.userId);
    
    this.chatService.chats$.subscribe(chats => {
      this.chatList = chats;
    });

    this.activatedRoute.params.subscribe(params => {
      const chatId = +params['id'];
      if (chatId) {
        this.chatService.readMessages(chatId);
        this.loadChatData(chatId);
        this.websocketService.connect(chatId);
      }
    });

    this.refreshInterval = setInterval(() => { // обновляем чаты раз в 1 секунду
      this.chatService.getAllChats(this.userId);
    }, 1000);
  }

  ngOnDestroy(): void {
    this.websocketService.closeConnection();
    if (this.refreshInterval) {
      clearInterval(this.refreshInterval);
    }
  }
  
  // Отправляем сообщение через вебсокет
  onSubmit(event?: Event) {
    if (event) {
      event.preventDefault(); // Предотвращает добавление символа "Enter"
    }
    this.messageText = this.sendMsgForm.controls['msgText'].value;
    this.sendMsgForm.reset();
    if (this.messageText && this.messageText.trim()) {
      const msgToSend = { 
        message_content: this.messageText,
        user_id: this.userId,
        chat_id: this.currentChat.id,
        } as Message;
      this.websocketService.sendMessage(JSON.stringify([msgToSend]));
      this.messageText = '';
    }
  }

  formatTime(createdAt: Date): string {
    return new Date(createdAt).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }

  selectedUserIds: number[] = [];

  onUserSelectionChanged(userId: number) {
    if (userId > 0) {
      this.selectedUserIds.push(userId); // Добавляем ID пользователя
    } else {
      const index = this.selectedUserIds.indexOf(-userId);
      if (index > -1) {
        this.selectedUserIds.splice(index, 1); // Удаляем ID пользователя
      }
    }
  }

  notCreatingChat() {
    this.createChat = false;
    this.selectedUserIds = [];
    this.countInArray.update(() => 0);
  }

  selectChat(chat: Chat) {
    this.currentChat = chat;
    this.router.navigateByUrl(`/chats/${chat.id}`);
    this.notCreatingChat();
  }

  isActiveChat(chat: Chat): boolean {
    return this.currentChat.id === chat.id;
  }

  onCreateChat() {
    if (this.selectedUserIds.length > 0) {
      this.chatService.createChat(this.selectedUserIds, this.userId).subscribe(
        (response) => {
          if (response.id) {
            this.websocketService.closeConnection();
            this.router.navigateByUrl(`/chats/${response.id}`).then(() => {
              window.location.reload();
            });
          }
        },
      );
    }
  }

  loadChatData(chatId: number): void {
    this.websocketService.closeConnection();
    this.chatService.getChatById(chatId).subscribe({
      next: (chat: Chat) => {
        this.selectChat({
          ...chat, 
          title: this.chatService.getTitle(chat.users, this.userId),
          });
        this.loadMessages(chatId);
      },
      error: (err) => {
        console.error('Ошибка при загрузке данных чата:', err);
      }
    });
  }
  
  loadMessages(chatId: number): void {
    this.chatService.getMessagesByChatId(chatId).subscribe({
      next: (messages: Message[]) => {
        this.messages = messages;
      },
      error: (err) => {
        console.error('Ошибка при загрузке сообщений:', err);
      }
    });
  }
}