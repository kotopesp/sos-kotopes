import {Component, Input, LOCALE_ID} from '@angular/core';
import { FavoritesButtonComponent } from "./favorites-button/favorites-button.component";
import { FindStatusFlagComponent } from "./find-status-flag/find-status-flag.component";
import {Post} from "../../model/post.interface";
import {UserService} from "../../services/user-service/user.service";
import {DatePipe, registerLocaleData} from "@angular/common";
import localeRu from '@angular/common/locales/ru';

registerLocaleData(localeRu);

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [FavoritesButtonComponent, FindStatusFlagComponent, DatePipe],
  templateUrl: './post.component.html',
  styleUrl: './post.component.scss',
  providers: [{ provide: LOCALE_ID, useValue: 'ru' }],
})
export class PostComponent {
  @Input() post!: Post;

  constructor(private userService: UserService) {
  }

}
