import {Routes} from '@angular/router';
import {PostsComponent} from "../pages/posts/ui/posts.component";
import {StartPageComponent} from "../pages/start-page/ui/start-page.component";
import {ChatTypeButtonComponent} from "../entities/chat-type-button/chat-type-button.component";
import {TestPageComponent} from "../pages/test-page/test-page.component";

export const routes: Routes = [
  {
    path: '', component: StartPageComponent
  },
  {
    path: 'posts', component: PostsComponent
  },
  {
    path: 'test', component: TestPageComponent
  }
];
