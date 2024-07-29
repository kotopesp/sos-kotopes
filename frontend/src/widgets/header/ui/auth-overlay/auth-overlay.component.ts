import { Component } from '@angular/core';
import {NgIf} from "@angular/common";

@Component({
  selector: 'app-auth-overlay',
  standalone: true,
  imports: [
    NgIf
  ],
  templateUrl: './auth-overlay.component.html',
  styleUrl: './auth-overlay.component.scss'
})
export class AuthOverlayComponent {
  passwordValid: boolean = true;
}
