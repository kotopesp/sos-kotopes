import {Component, inject, Input, signal, WritableSignal} from '@angular/core';
import {NgForOf} from "@angular/common";
import {AuthService} from "../../../../services/auth-service/auth.service";
import {UserService} from "../../../../services/user-service/user.service";

@Component({
  selector: 'app-profile-popup',
  standalone: true,
  imports: [
    NgForOf
  ],
  templateUrl: './profile-popup.component.html',
  styleUrl: './profile-popup.component.scss'
})
export class ProfilePopupComponent {
  @Input() isAuth: WritableSignal<boolean>;

  constructor(private authService: AuthService, private userService: UserService) {
    this.isAuth = signal<boolean>(false);
  }

  controlItems = [
    {
      path: '',
      title: 'Избранное',
      icon: 'heart-icon.svg',
    },
    {
      path: '',
      title: 'Профиль',
      icon: 'profile-icon.svg',
    },
    {
      path: '',
      title: 'Настройки',
      icon: 'settings-icon.svg',
    },
    {
      path: '',
      title: 'Выйти',
      icon: 'logout-icon.svg',
    },
  ]

  logout() {
    this.authService.logout();
    this.isAuth.set(false);
  }
}
