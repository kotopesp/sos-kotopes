import { Component } from '@angular/core';
import {ChooseOneDirective} from "../../../directives/choose-one/choose-one.directive";

@Component({
  selector: 'app-button-lost-pet',
  standalone: true,
    imports: [
        ChooseOneDirective
    ],
  templateUrl: './button-lost-pet.component.html',
  styleUrl: './button-lost-pet.component.scss'
})
export class ButtonLostPetComponent {

}
