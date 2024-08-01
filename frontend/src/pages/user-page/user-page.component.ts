import { Component } from '@angular/core';

import {HeaderComponent} from "../../widgets/header/header.component";
import { UserService } from '../../services/user-service/user.service';
import { RoleButtonComponent } from './role-button/role-button.component';
import { IconButtonComponent } from './icon-button/icon-button.component';

@Component({
  selector: 'app-user-page',
  standalone: true,
  imports: [HeaderComponent, RoleButtonComponent, IconButtonComponent],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss'
})
export class UserPageComponent {
  firstName = 'Тимофей'
  secondName = 'Зайнулин'
  onlineStatus = 'В сети'
  username = 'tim.violine'
  totalPosts = '22'
  profilePhoto = '../../assets/images/test-cat.png'

  isOwnAccount = false
}
