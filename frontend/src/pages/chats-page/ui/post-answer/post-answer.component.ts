import { Component } from '@angular/core';
import {
  FindStatusButtonsComponent
} from "../../../../widgets/filters-bar/find-status/find-status-buttons/find-status-buttons.component";
import {PostStatusComponent} from "../../../../shared/post-status/post-status.component";

@Component({
  selector: 'app-post-answer',
  standalone: true,
  imports: [
    FindStatusButtonsComponent,
    PostStatusComponent
  ],
  templateUrl: './post-answer.component.html',
  styleUrl: './post-answer.component.scss'
})
export class PostAnswerComponent {

}
