import {Routes} from '@angular/router';
import {AdsComponent} from "./ads/ads.component";
import {StartPageComponent} from "./start-page/start-page.component";

export const routes: Routes = [
  {
    path: '', component: StartPageComponent
  },
  {
    path: 'ads', component: AdsComponent
  }
];
