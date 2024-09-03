import {Component, OnInit, signal} from '@angular/core';
import { ActivatedRoute, Params } from '@angular/router';
import {HeaderComponent} from "../../widgets/header/header.component";

import { RoleButtonComponent } from './role-button/role-button.component';

import { map, Observable, switchMap } from 'rxjs';
import { CommonModule } from '@angular/common';
import {UserService} from "../../services/user.service";
import {User} from "../../model/user.interface";
import {PostComponent} from "../../shared/post/post.component";
import {WriteOverlayComponent} from "./write-overlay/write-overlay.component";

@Component({
  selector: 'app-user-page',
  standalone: true,
  imports: [HeaderComponent, RoleButtonComponent, CommonModule, PostComponent, WriteOverlayComponent],
  templateUrl: './user-page.component.html',
  styleUrl: './user-page.component.scss'
})
export class UserPageComponent implements OnInit {
  firstName = 'Тимофей';
  secondName = 'Зайнулин';
  onlineStatus = 'В сети';
  username = 'tim.violine';
  totalPosts = '22';
  profilePhoto = '../../assets/images/test-cat.png';
  isOwnAccount = false;

  user$!: Observable<User>;

  likeActive = false;
  editTextArea = false;
  pressedPost  = true;
  haveRoles = true;
  openInfo: number | null = null;

  writeOverlay = signal<boolean>(false);

  constructor(private activatedRoute: ActivatedRoute, private userService: UserService) {
  }

  ngOnInit(): void {
    this.user$ = this.activatedRoute.params.pipe(
      map((params: Params) => parseInt(params['id'], 10)),
      switchMap((userId: number) => this.userService.getById(userId))
    );
  }

  likeActiveButton(): void {
    this.likeActive = !this.likeActive;
  }

  onClickFirst(): void {
    if (this.openInfo === 1) {
      this.openInfo = null;
    } else {
      this.openInfo = 1;
    }
  }
  onClickSecond() {
    if (this.openInfo === 2) {
      this.openInfo = null;
    } else {
      this.openInfo = 2;
    }
  }
  onClickThird() {
    if (this.openInfo === 3) {
      this.openInfo = null;
    } else {
      this.openInfo = 3;
    }
  }
}
