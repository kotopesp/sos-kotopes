import {Component, Input, signal, WritableSignal} from '@angular/core';
import {NgIf} from "@angular/common";
import {FormControl, FormGroup, ReactiveFormsModule, Validators} from "@angular/forms";
import {RegisterService} from "../../../../services/register-service/register.service";
import {Router} from "@angular/router";


@Component({
  selector: 'app-register-overlay',
  standalone: true,
  imports: [
    NgIf,
    ReactiveFormsModule
  ],
  templateUrl: './register-overlay.component.html',
  styleUrl: './register-overlay.component.scss'
})
export class RegisterOverlayComponent {
  @Input() isRegisterOverlay: WritableSignal<boolean>;
  @Input() isAuthOverlay: WritableSignal<boolean>;
  formRegister: FormGroup;
  isPasswordVisible = signal<boolean>(false);
  isPasswordVisibleRepeat = signal<boolean>(false);

  constructor(private register: RegisterService, private router: Router) {
    this.isRegisterOverlay = signal<boolean>(false)
    this.isAuthOverlay = signal<boolean>(false)

    this.formRegister = new FormGroup({
      email: new FormControl(null, Validators.required),
      name: new FormControl(null, Validators.required),
      lastname: new FormControl(null, Validators.required),
      username: new FormControl(null, Validators.required),
      password: new FormControl(null, Validators.required),
    })
  }

  onSubmit() {

    this.register.registration(this.formRegister.value).subscribe(
      {
        next: () => {
          this.isRegisterOverlay.set(false);
        },
        error: (error) => {
          console.warn(error);
        },
     }
    )
  }
}

