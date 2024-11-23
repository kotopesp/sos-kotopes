import {Component} from '@angular/core';
import { PostComponent } from "../../../shared/post/post.component";
import {FiltersBarComponent} from "../../../widgets/filters-bar/filters-bar.component";
import {RouterLink} from "@angular/router";
import {PostsService} from "../../../services/posts-services/posts.service";
import {Meta, Post} from "../../../model/post.interface";
import {AsyncPipe, NgForOf, NgIf} from "@angular/common";

@Component({
  selector: 'app-posts-services',
  standalone: true,
  imports: [
    PostComponent,
    FiltersBarComponent,
    RouterLink,
    NgForOf,
    AsyncPipe,
    NgIf
  ],
  templateUrl: './posts.component.html',
  styleUrl: './posts.component.scss'
})
export class PostsComponent {
  posts: Post[] = [];
  meta: Meta | null = null;
  constructor(private postService: PostsService) {
    this.postService.getPosts().subscribe(response => {
      if (response) {
        this.posts = response.data.posts;  // Сохраняем массив постов
        this.meta = response.data.meta;
      }
    });
  }

}
