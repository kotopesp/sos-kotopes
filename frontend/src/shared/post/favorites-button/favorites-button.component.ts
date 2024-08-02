import { Component } from '@angular/core';
import {NgClass} from "@angular/common";

@Component({
  selector: 'favorites-button',
  standalone: true,
  imports: [NgClass],
  templateUrl: './favorites-button.component.html',
  styleUrl: './favorites-button.component.scss'
})
export class FavoritesButtonComponent {
  isPressed: boolean = false;

  pressButton() {
    this.isPressed = !this.isPressed;
  }
}
