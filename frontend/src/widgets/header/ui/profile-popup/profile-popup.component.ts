import {Component, Input, signal, WritableSignal} from '@angular/core';
import {NgForOf} from "@angular/common";
import {AuthService} from "../../../../services/auth-service/auth.service";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-profile-popup',
  standalone: true,
  imports: [
    NgForOf,
    RouterLink
  ],
  templateUrl: './profile-popup.component.html',
  styleUrl: './profile-popup.component.scss'
})
export class ProfilePopupComponent {
  @Input() isAuth: WritableSignal<boolean>;

  constructor(private authService: AuthService) {
    this.isAuth = signal<boolean>(false);
  }

  controlItems = [
    {
      path: '',
      title: 'Избранное',
      icon: 'heart-icon.svg',
    },
    {
      path: 'users/1',
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
