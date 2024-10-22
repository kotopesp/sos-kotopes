import { Component } from '@angular/core';
import {NgForOf} from "@angular/common";

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
  controlItems = [
    {
      path: '',
      title: 'Избранное',
      icon: 'heart-icon.svg'
    },
    {
      path: '',
      title: 'Профиль',
      icon: 'profile-icon.svg'
    },
    {
      path: '',
      title: 'Настройки',
      icon: 'settings-icon.svg'
    },
    {
      path: '',
      title: 'Выйти',
      icon: 'logout-icon.svg'
    },
  ]
}
