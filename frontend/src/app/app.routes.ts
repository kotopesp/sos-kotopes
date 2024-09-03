import {Routes} from '@angular/router';

import { PostsComponent } from "../pages/posts/ui/posts.component";
import { StartPageComponent } from "../pages/start-page/ui/start-page.component";
import { UserPageComponent } from '../pages/user-page/user-page.component';
import {ChatsPageComponent} from "../pages/chats-page/chats-page.component";
export const routes: Routes = [
  {
    path: '', component: StartPageComponent
  },
  {
    path: 'posts', component: PostsComponent
  },
  {
    path: 'chats', component: ChatsPageComponent
  },
  {
    path: 'users/:id', component: UserPageComponent
  }
];
