import {Component} from '@angular/core';
import { FavoritesButtonComponent } from "./favorites-button/favorites-button.component";
import { FindStatusComponent } from "../../../../widgets/filters-bar/find-status/find-status.component";

@Component({
  selector: 'app-post',
  standalone: true,
  imports: [FavoritesButtonComponent, FindStatusComponent],
  templateUrl: './post.component.html',
  styleUrl: './post.component.scss'
})
export class PostComponent {

}
