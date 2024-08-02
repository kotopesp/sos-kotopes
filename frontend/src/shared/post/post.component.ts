import { Component } from '@angular/core';
import { FavoritesButtonComponent } from "./favorites-button/favorites-button.component";
import { FindStatusFlagComponent } from "./find-status-flag/find-status-flag.component";

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [FavoritesButtonComponent, FindStatusFlagComponent],
  templateUrl: './post.component.html',
  styleUrl: './post.component.scss'
})
export class PostComponent {

}
