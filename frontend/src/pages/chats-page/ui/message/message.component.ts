import {Component, Input} from '@angular/core';
import { Timestamp } from 'rxjs';

@Component({
  selector: 'app-message',
  standalone: true,
  imports: [],
  templateUrl: './message.component.html',
  styleUrl: './message.component.scss'
})
export class MessageComponent {
  @Input() answer: boolean = false;
  @Input() messageContent: string = '';
  @Input() messageTime: string = '';
}
