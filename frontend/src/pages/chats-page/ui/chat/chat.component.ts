import {Component} from '@angular/core';
import {NgClass} from "@angular/common";
import {ToggleActiveDirective} from "../../../../directives/toggle-active/toggle-active.directive";

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
export class ChatComponent {
}
