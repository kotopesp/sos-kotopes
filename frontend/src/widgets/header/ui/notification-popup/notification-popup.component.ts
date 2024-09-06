import { Component } from '@angular/core';
import {NotificationRespondedComponent} from "./ui/notification-responded/notification-responded.component";
import {NotificationAddFavoritesComponent} from "./ui/notification-add-favorites/notification-add-favorites.component";
import {NotificationLeftCommentComponent} from "./ui/notification-left-comment/notification-left-comment.component";
import {NotificationPostClosedComponent} from "./ui/notification-post-closed/notification-post-closed.component";

@Component({
  selector: 'app-notification-popup',
  standalone: true,
  imports: [
    NotificationRespondedComponent,
    NotificationAddFavoritesComponent,
    NotificationLeftCommentComponent,
    NotificationPostClosedComponent
  ],
  templateUrl: './notification-popup.component.html',
  styleUrl: './notification-popup.component.scss'
})
export class NotificationPopupComponent {

}
