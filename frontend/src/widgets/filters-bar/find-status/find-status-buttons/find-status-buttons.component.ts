import { Component } from '@angular/core';
import { NgClass } from "@angular/common";

@Component({
  selector: 'app-find-status-botton',
  standalone: true,
  imports: [NgClass],
  templateUrl: './find-status-botton.component.html',
  styleUrl: './find-status-botton.component.scss'
})
export class FindStatusBottonComponent {
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
