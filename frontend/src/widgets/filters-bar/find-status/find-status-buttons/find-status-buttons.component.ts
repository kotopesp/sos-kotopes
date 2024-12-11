import { Component } from '@angular/core';
import { NgClass } from "@angular/common";

@Component({
  selector: 'app-find-status-buttons',
  standalone: true,
  imports: [NgClass],
  templateUrl: './find-status-buttons.component.html',
  styleUrl: './find-status-buttons.component.scss'
})
export class FindStatusButtonsComponent {
  isPressedLost = false
  pressButtonLost() {
    this.isPressedLost = !this.isPressedLost
  }

  isPressedFoundHome = false
  pressButtonFoundHome() {
    this.isPressedFoundHome = !this.isPressedFoundHome
  }

  isPressedLookingFor = false
  pressButtonLookingFor() {
    this.isPressedLookingFor = !this.isPressedLookingFor
  }
}
