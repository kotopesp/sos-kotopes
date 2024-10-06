import { Component } from '@angular/core';
import {ChooseOneDirective} from "../../../directives/choose-one/choose-one.directive";

@Component({
  selector: 'app-button-looking-for-home',
  standalone: true,
    imports: [
        ChooseOneDirective
    ],
  templateUrl: './button-looking-for-home.component.html',
  styleUrl: './button-looking-for-home.component.scss'
})
export class ButtonLookingForHomeComponent {

}
