import { Component } from '@angular/core';
import { NgClass } from "@angular/common";

@Component({
  selector: 'find-status-buttons',
  standalone: true,
  imports: [NgClass],
  templateUrl: './find-status-buttons.component.html',
  styleUrl: './find-status-buttons.component.scss'
})
export class FindStatusButtonsComponent {
  isPressedLost: boolean = false
  pressButtonLost() {
    this.isPressedLost = !this.isPressedLost
  }

  isPressedFoundHome:boolean = false
  pressButtonFoundHome() {
    this.isPressedFoundHome = !this.isPressedFoundHome
  }

  isPressedLookingFor: boolean = false
  pressButtonLookingFor() {
    this.isPressedLookingFor = !this.isPressedLookingFor
  }
}
