import {AfterViewChecked, Component, ElementRef, signal, ViewChild} from '@angular/core';
import {AppChatTypeButtonComponent} from "../../entities/chat-type-button/app-chat-type-button.component";
import {Button} from "../../model/button";
import {NgClass, NgForOf, NgIf, NgStyle, NgSwitch, NgSwitchCase} from "@angular/common";
import {ChatComponent} from "./ui/chat/chat.component";
import {MessageComponent} from "./ui/message/message.component";
import {PostAnswerComponent} from "./ui/post-answer/post-answer.component";
import {AddToChatComponent} from "./ui/add-to-chat/add-to-chat.component";
import {ToggleActiveDirective} from "../../directives/toggle-active/toggle-active.directive";

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
  ],
  templateUrl: './chats-page.component.html',
  styleUrl: './chats-page.component.scss'
})
export class ChatsPageComponent implements AfterViewChecked {
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

  ngAfterViewChecked(): void {
    this.scrollToBottom();
  }

  private scrollToBottom(): void {
    if (this.scrollableContainer && this.scrollableContainer.nativeElement) {
      const container = this.scrollableContainer.nativeElement as HTMLElement;
      container.scrollTop = container.scrollHeight;
    }
  }
}
