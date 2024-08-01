import { Component } from '@angular/core';

import {HeaderComponent} from "../../widgets/header/header.component";
import { UserService } from '../../services/user-service/user.service';
import { RoleButtonComponent } from './role-button/role-button.component';

@Component({
  selector: 'app-user-page',
  standalone: true,
  imports: [HeaderComponent, RoleButtonComponent],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss'
})
export class UserPageComponent {
  firstName = 'John'
  secondName = 'Johnson'
  onlineStatus = 'online'
  username = 'tim.violine'
  totalPosts = '22'
  profilePhoto = '../../assets/images/test-cat.png'
}
