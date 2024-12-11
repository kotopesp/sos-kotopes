import {Component, Input} from "@angular/core";
import {ToggleActiveDirective} from "../../pages/chats-page/toggle-active.directive";

@Component({
  selector: 'app-chat-type-button',
  standalone: true,
  templateUrl: "app-chat-type-button.component.html",
  imports: [
    ToggleActiveDirective
  ],
  styles: `
    .chat-type-icon {
      margin-right: 20px;
    }

    .counter {
      display: flex;
      align-items: center;
      justify-content: center;
      width: 30px;
      height: 30px;
      background: white;
      opacity: 70%;
      font-size: 18px;
      border-radius: 20px;
      margin-right: 15px;
    }

    .chat-type-button {
      display: flex;
      align-items: center;
      height: 100px;
      width: 100%;
      font-size: 20px;
      border-radius: 100px;
      padding-left: 30px;
      margin-bottom: 10px;
    }
  `
})
export class AppChatTypeButtonComponent {
  @Input({required: true})
  label = ''

  @Input({required: true})
  color = ''

  @Input()
  counter = 0

  @Input()
  icon = "allChats.svg"
}
