import {Component, Input, signal, WritableSignal} from '@angular/core';
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
    ReactiveFormsModule,
  ],
  templateUrl: './auth-service-overlay.component.html',
  styleUrl: './auth-service-overlay.component.scss'
})
export class AuthServiceOverlayComponent {
  @Input() isAuthOverlay: WritableSignal<boolean>;
  @Input() isRegisterOverlay: WritableSignal<boolean>;
  @Input() isAuth: WritableSignal<boolean>;
  passwordValid  = true;
  isPasswordVisible: WritableSignal<boolean> = signal<boolean>(false);
  formAuth: FormGroup;


  constructor(private auth: AuthService) {
    this.isAuthOverlay = signal<boolean>(false);
    this.isRegisterOverlay = signal<boolean>(false);
    this.isAuth = signal<boolean>(false);

    this.formAuth = new FormGroup({
      username: new FormControl(null, Validators.required),
      password: new FormControl(null, Validators.required)
    })
  }

  onSubmit() {
    this.formAuth.disable()
    this.auth.login(this.formAuth.value).subscribe(
      {
        next: () => {
          this.isAuthOverlay.set(false);
          this.isAuth.set(true);
          },
        error: (error) => {
          console.warn(error);
          this.formAuth.enable()
        }
      }
    )
  }
}
