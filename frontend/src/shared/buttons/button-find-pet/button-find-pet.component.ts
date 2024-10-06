import { Component } from '@angular/core';
import {ChooseOneDirective} from "../../../directives/choose-one/choose-one.directive";

@Component({
  selector: 'app-button-find-pet',
  standalone: true,
  imports: [
    ChooseOneDirective
  ],
  templateUrl: './button-find-pet.component.html',
  styleUrl: './button-find-pet.component.scss'
})
export class ButtonFindPetComponent {

}
