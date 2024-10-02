import {Component, Input} from '@angular/core';
import { NgIf } from '@angular/common';
@Component({
  selector: 'app-message',
  standalone: true,
  imports: [NgIf],
  templateUrl: './message.component.html',
  styleUrl: './message.component.scss'
})
export class MessageComponent {
  @Input() answer: boolean = false;
  @Input() messageContent: string = '';
  @Input() messageTime: string = '';
  @Input() name: string = '';
}
