import {Component} from '@angular/core';
import {HeaderComponent} from "../../../widgets/header/header.component";
import {PostComponent} from "./post/post.component";

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [
    HeaderComponent,
    PostComponent
  ],
  templateUrl: './posts.component.html',
  styleUrl: './posts.component.scss'
})
export class PostsComponent {

}
