import { Component } from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import {HeaderComponent} from "../../widgets/header/header.component";
import { UserService } from '../../services/user-service/user.service';
import { RoleButtonComponent } from './role-button/role-button.component';
import { IconButtonComponent } from './icon-button/icon-button.component';
import { User } from '../../model/user.interface';
import { map, Observable, switchMap } from 'rxjs';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-user-page',
  standalone: true,
  imports: [HeaderComponent, RoleButtonComponent, IconButtonComponent, CommonModule],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss'
})
export class UserPageComponent {
  firstName = 'Тимофей';
  secondName = 'Зайнулин';
  onlineStatus = 'В сети';
  username = 'tim.violine';
  totalPosts = '22';
  profilePhoto = '../../assets/images/test-cat.png';
  isOwnAccount = false;

  user$!: Observable<User>;

  constructor(private activatedRoute: ActivatedRoute, private userService: UserService) { }

  ngOnInit(): void {
    this.user$ = this.activatedRoute.params.pipe(
      map((params: Params) => parseInt(params['id'], 10)),
      switchMap((userId: number) => this.userService.getById(userId))
    );
  }
}
