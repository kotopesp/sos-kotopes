import {Component} from '@angular/core';
import {HeaderComponent} from "../../../widgets/header/header.component";
import { PostComponent } from "../../../shared/post/post.component";
import {FiltersBarComponent} from "../../../widgets/filters-bar/filters-bar.component";
import {RouterLink} from "@angular/router";

@Component({
  selector: 'app-posts-services',
  standalone: true,
  imports: [
    HeaderComponent,
    PostComponent,
    FiltersBarComponent,
    RouterLink
  ],
  templateUrl: './posts.component.html',
  styleUrl: './posts.component.scss'
})
export class PostsComponent {

}
