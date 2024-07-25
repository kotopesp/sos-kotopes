import {Component} from '@angular/core';
import {HeaderComponent} from "../../../widgets/header/header.component";
import {PostComponent} from "./post/post.component";
import { SidebarComponent } from "../../../widgets/sidebar/sidebar.component";

@Component({
  selector: 'app-posts',
  standalone: true,
  imports: [
    HeaderComponent,
    PostComponent,
    SidebarComponent
],
  templateUrl: './posts.component.html',
  styleUrl: './posts.component.scss'
})
export class PostsComponent {

}
