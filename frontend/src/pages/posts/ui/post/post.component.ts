import {Component} from '@angular/core';
import { FavoritesButtonComponent } from "./favorites-button/favorites-button.component";

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [FavoritesButtonComponent],
  templateUrl: './post.component.html',
  styleUrl: './post.component.scss'
})
export class PostComponent {

}
