import {Component, inject, Input, signal, WritableSignal} from '@angular/core';
import {NgIf} from "@angular/common";
import {RegisterOverlayComponent} from "../register-overlay/register-overlay.component";
import {AuthService} from "../../../../services/auth-service/auth.service";
import {FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";

@Component({
  selector: 'app-auth-service-overlay',
  standalone: true,
  imports: [
    NgIf,
    RegisterOverlayComponent,
    ReactiveFormsModule
  ],
  templateUrl: './auth-service-overlay.component.html',
  styleUrl: './auth-service-overlay.component.scss'
})
export class AuthServiceOverlayComponent {
  @Input() isAuthOverlay: WritableSignal<boolean>;
  @Input() isRegisterOverlay: WritableSignal<boolean>;
  passwordValid  = true;
  isPasswordVisible: WritableSignal<boolean> = signal<boolean>(false);

  formAuth = new FormGroup({
    username: new FormControl(null, Validators.required),
    password: new FormControl(null, Validators.required)
  })

  authService = inject(AuthService)

  constructor() {
    this.isAuthOverlay = signal<boolean>(false);
    this.isRegisterOverlay = signal<boolean>(false);
  }

  onSubmit() {

    if (this.formAuth.valid) {

      this.authService.login(this.formAuth.value)
    }
  }
}
