import {Component, signal} from '@angular/core';
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-register-overlay',
  standalone: true,
  imports: [
    NgIf
  ],
  templateUrl: './register-overlay.component.html',
  styleUrl: './register-overlay.component.scss'
})
export class RegisterOverlayComponent {
  isPasswordVisible = signal<boolean>(false);
  isPasswordVisibleRepeat = signal<boolean>(false);
}
