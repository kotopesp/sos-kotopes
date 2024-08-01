import { Component } from '@angular/core';

import {HeaderComponent} from "../../widgets/header/header.component";
import { UserService } from '../../services/user-service/user.service';

@Component({
  selector: 'app-user-page',
  standalone: true,
  imports: [HeaderComponent],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss'
})
export class UserPageComponent {
  
}
