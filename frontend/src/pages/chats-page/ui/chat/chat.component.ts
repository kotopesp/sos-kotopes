import {Component, Input} from '@angular/core';
import {NgClass, NgIf} from "@angular/common";
import {ToggleActiveDirective} from "../../toggle-active.directive";
import { Chat } from '../../../../model/chat.interface';

@Component({
  selector: 'app-chat',
  standalone: true,
  imports: [
    NgClass,
    NgIf,
    ToggleActiveDirective,
  ],
  templateUrl: './chat.component.html',
  styleUrl: './chat.component.scss'
})
export class ChatComponent {
  @Input() chat!: Chat;  
}