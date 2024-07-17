import { Component } from '@angular/core';
import {HeaderComponent} from "../common-ui/header/header.component";
import {AdPetComponent} from "../common-ui/ad-pet/ad-pet.component";

@Component({
  selector: 'app-ads',
  standalone: true,
  imports: [
    HeaderComponent,
    AdPetComponent
  ],
  templateUrl: './ads.component.html',
  styleUrl: './ads.component.scss'
})
export class AdsComponent {

}
