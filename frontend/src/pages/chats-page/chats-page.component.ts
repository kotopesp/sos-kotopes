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
import { map, Observable, switchMap } from 'rxjs';
import { User } from '../../model/user.interface';
import { ActivatedRoute, Params, Router } from '@angular/router';
import {ChatService} from '../../services/chat-service/chat.service';
import { HttpClient } from '@angular/common/http';
import { Chat } from '../../model/chat.interface';
import { Chatinfo } from '../../model/chatinfo.interface';

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
  currentChat: Chatinfo = { Id: -1, Title: '', Chattype: '' };
  activeChatId: number = -1;
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

  ngAfterViewChecked(): void {
    this.scrollToBottom();
  }

  private scrollToBottom(): void {
    if (this.scrollableContainer && this.scrollableContainer.nativeElement) {
      const container = this.scrollableContainer.nativeElement as HTMLElement;
      container.scrollTop = container.scrollHeight;
    }
  }

  public messages: { content: string, isUserMessage: boolean, time: string }[] = [];
  public messageText: string = '';
  public userId: string = '1'; // id пользователя
  public user$!: Observable<User>;
  public favusers: { id: number, username: string}[] = [];
  public chatList: Chatinfo[] = [];

  constructor(private router: Router, private activatedRoute: ActivatedRoute, private chatService: ChatService, private websocketService: WebsocketService, private http: HttpClient) {
    // this.user$.subscribe((user: User) => {
    //   if (user && user.id) {
    //     this.userId = user.id.toString();  // Сохраняем id пользователя как строку
    //   }
    // });
  }

  ngOnInit(): void {
    this.user$ = this.activatedRoute.params.pipe(
      map((params: Params) => parseInt(params['id'], 10)),
      switchMap((chatId: number) => this.chatService.getById(chatId))
    );

    // Получаем входящие сообщения
    this.websocketService.getMessages().subscribe((msg: string) => {
      const parsedMsg = this.parseMessage(msg);
      this.messages.push(parsedMsg);
    });

    this.chatService.getFavUsers().subscribe(
      (users) => {
        this.favusers = users; // Заполняем массив пользователей
      }
    );

    this.chatService.getAllChats().subscribe(
      (chats) => {
        this.chatList = chats;
      }
    );

    this.activatedRoute.params.subscribe(params => {
      const chatId = +params['id']; // Получаем ID чата из маршрута
      if (chatId) {
        this.loadChatData(chatId); // Загружаем данные чата
      }
    });
  }
  
  // Отправляем сообщение через вебсокет
  onSubmit() {
    this.messageText = this.sendMsgForm.controls['msgText'].value;
    this.sendMsgForm.reset();
    if (this.messageText.trim()) {
      const timeNow = Date.now();
      this.websocketService.sendMessage(JSON.stringify([
        { content: this.messageText, userId: this.userId, createdAt: timeNow, updatedAt: timeNow }, // Message form
      ]));
      this.messageText = '';
    }
  }

  private parseMessage(msg: string): { content: string, isUserMessage: boolean, time: string } {
    const msgJson = JSON.parse(msg);
    // console.log("msg id: ", msgJson[0].userId, "cur usr id: ", this.userId, "eqls: ", msgJson[0].userId === this.userId);
    var msgDate = new Date(msgJson[0].updatedAt);
    var timeString = `${msgDate.getHours()}:${msgDate.getMinutes()}`
    return {
      content: msgJson[0].content,
      isUserMessage: msgJson[0].userId === this.userId,
      time: timeString,
    };
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
    this.countInArray.update(_ => 0);
  }

  selectChat(chat: Chatinfo) {
    this.currentChat = chat;
    this.activeChatId = chat.Id;
    this.notCreatingChat();
  }

  isActiveChat(chat: Chatinfo): boolean {
    // console.log(this.activeChatId, "cur chat:", chat.Id)
    return this.activeChatId === chat.Id;
  }

  onCreateChat() {
    if (this.selectedUserIds.length > 0) {
      this.chatService.createChat(this.selectedUserIds).subscribe(
        (response) => {
          if (response.Id) {
            // console.log("перенаправляем на", response);
            this.router.navigateByUrl(`/chats/${response.Id}`).then(() => {
              window.location.reload();
            });
          }
        },
      );
    }
    else {
      console.log("no selected users");
    }
  }

  loadChatData(chatId: number): void {
    // Получение информации о чате через сервис чатов
    this.chatService.getChatById(chatId).subscribe({
      next: (chat: Chat) => {
        this.selectChat({Id: chat.ID, Chattype: chat.ChatType, Title: this.chatService.getTitle(chat.Users)});

        // this.loadMessages(chatId);
      },
      error: (err) => {
        console.error('Ошибка при загрузке данных чата:', err);
      }
    });
  }
  
  // loadMessages(chatId: number): void {
  //   this.chatService.getMessagesByChatId(chatId).subscribe({
  //     next: (messages: Message[]) => {
  //       this.messages = messages;

  //       this.scrollToBottom();
  //     },
  //     error: (err) => {
  //       console.error('Ошибка при загрузке сообщений:', err);
  //     }
  //   });
  // }
}