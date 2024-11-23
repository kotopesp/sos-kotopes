import {Component, Input} from '@angular/core';
import {NgClass} from "@angular/common";
import {PostsService} from "../../../services/posts-services/posts.service";

@Component({
  selector: 'app-favorites-button',
  standalone: true,
  imports: [NgClass],
  templateUrl: './favorites-button.component.html',
  styleUrl: './favorites-button.component.scss'
})
export class FavoritesButtonComponent {
  @Input() isPressed!: boolean;
  @Input() postID: number = 0;

  constructor(private postService: PostsService) {
  }

  pressButton() {
    this.isPressed = !this.isPressed;
    if (this.isPressed) {
      // add to favorites
      this.postService.addPostToFavorites(1)
    } else {
      // delete from favorites
      this.postService.deletePostFromFavorites(1)
    }
  }
}
