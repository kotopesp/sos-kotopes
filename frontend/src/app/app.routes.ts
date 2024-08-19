import {Routes} from '@angular/router';
<<<<<<< HEAD
import {PostsComponent} from "../pages/posts/ui/posts.component";
import {StartPageComponent} from "../pages/start-page/ui/start-page.component";
import { UserPageComponent } from '../pages/user-page/ui/user-page.component';
=======
import { PostsComponent } from "../pages/posts/ui/posts.component";
import { StartPageComponent } from "../pages/start-page/ui/start-page.component";
import { UserPageComponent } from '../pages/user-page/user-page.component';
>>>>>>> origin/frontend

export const routes: Routes = [
  {
    path: '', component: StartPageComponent
  },
  {
    path: 'posts', component: PostsComponent
  },
  {
    path: 'users/:id', component: UserPageComponent
  }
];
