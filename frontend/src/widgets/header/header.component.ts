import {Component, signal, WritableSignal} from '@angular/core';
import {RouterLink} from "@angular/router";
import {NgForOf, NgIf} from "@angular/common";
import {ProfilePopupComponent} from "./ui/profile-popup/profile-popup.component";
import {NotificationPopupComponent} from "./ui/notification-popup/notification-popup.component";
import {MessagePopupComponent} from "./ui/message-popup/message-popup.component";
import {AuthServiceOverlayComponent} from "./ui/auth-overlay/auth-service-overlay.component";
import {RegisterOverlayComponent} from "./ui/register-overlay/register-overlay.component";

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
    AuthServiceOverlayComponent,
    RegisterOverlayComponent,
  ],
  templateUrl: './header.component.html',
  styleUrl: './header.component.scss'
})
export class HeaderComponent {
  isAuth = false;
  isAuthOverlay: WritableSignal<boolean> = signal<boolean>(false);
  isRegisterOverlay: WritableSignal<boolean> = signal<boolean>(false);

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
      path: 'vets',
      text: 'Ветеринары',
      className: 'header__vets'
    },
    {
      path: '',
      text: 'Как я могу помочь?',
      className: 'header__how-to-help'
    },
  ]
}
