import {Component, Input, signal, WritableSignal} from '@angular/core';
import {NgIf} from "@angular/common";
import {RegisterOverlayComponent} from "../register-overlay/register-overlay.component";

@Component({
  selector: 'app-auth-overlay',
  standalone: true,
  imports: [
    NgIf,
    RegisterOverlayComponent
  ],
  templateUrl: './auth-overlay.component.html',
  styleUrl: './auth-overlay.component.scss'
})
export class AuthOverlayComponent {
  @Input() isAuthOverlay: any;
  passwordValid: boolean = true;
  isPasswordVisible: WritableSignal<boolean> = signal<boolean>(false);
}
