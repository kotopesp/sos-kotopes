import { Component } from '@angular/core';
import {NgClass} from "@angular/common";

@Component({
  selector: 'app-favorites-button',
  standalone: true,
  imports: [NgClass],
  templateUrl: './favorites-button.component.html',
  styleUrl: './favorites-button.component.scss'
})
export class FavoritesButtonComponent {
  isPressed = false;

  pressButton() {
    this.isPressed = !this.isPressed;
  }
}
