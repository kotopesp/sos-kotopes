import {Component} from '@angular/core';
import {RouterLink} from "@angular/router";
import {NgForOf, NgIf} from "@angular/common";
import {ProfilePopupComponent} from "./ui/profile-popup/profile-popup.component";
import {NotificationPopupComponent} from "./ui/notification-popup/notification-popup.component";
import {MessagePopupComponent} from "./ui/message-popup/message-popup.component";
import {AuthOverlayComponent} from "../auth-overlay/auth-overlay.component";

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [
    RouterLink,
    NgForOf,
    NgIf,
    ProfilePopupComponent,
    ProfilePopupComponent,
    NotificationPopupComponent,
    MessagePopupComponent,
    AuthOverlayComponent
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent {
  isAuth: boolean = true;

  headerItems = [
    {
      path: 'posts',
      text: 'Объявления',
      className: 'header__ads'
    },
    {
      path: '',
      text: 'Отловщики',
      className: 'header__overexposure'
    },
    {
      path: '',
      text: 'Передержка',
      className: 'header__ads'
    },
    {
      path: '',
      text: 'Как я могу помочь?',
      className: 'header__how-to-help'
    },
  ]
}
