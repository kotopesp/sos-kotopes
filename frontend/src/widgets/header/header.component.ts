import {Component, OnInit, signal, WritableSignal} from '@angular/core';
import {Router, RouterLink} from "@angular/router";
import {NgForOf, NgIf} from "@angular/common";
import {ProfilePopupComponent} from "./ui/profile-popup/profile-popup.component";
import {NotificationPopupComponent} from "./ui/notification-popup/notification-popup.component";
import {MessagePopupComponent} from "./ui/message-popup/message-popup.component";
import {AuthServiceOverlayComponent} from "./ui/auth-overlay/auth-service-overlay.component";
import {RegisterOverlayComponent} from "./ui/register-overlay/register-overlay.component";
import {AuthService} from "../../services/auth-service/auth.service";

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
export class HeaderComponent implements OnInit {
  isAuth: WritableSignal<boolean>;
  isAuthOverlay: WritableSignal<boolean>;
  isRegisterOverlay: WritableSignal<boolean>;

  constructor(private authService: AuthService, private router: Router) {
    this.isAuth = signal<boolean>(false);
    this.isAuthOverlay = signal<boolean>(false);
    this.isRegisterOverlay = signal<boolean>(false);
  }

  ngOnInit() {
    this.isAuth = signal<boolean>(this.authService.isAuth);
    this.isAuthOverlay = signal<boolean>(false);
    this.isRegisterOverlay = signal<boolean>(false);
  }

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

  redirectToChats(): void {
    this.router.navigate(['/chats']);
  }
}
