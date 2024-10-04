import {Component, Input, WritableSignal} from '@angular/core';
import {AddPhotoButtonComponent} from "../buttons/add-photo-button/add-photo-button.component";
import {NgIf, NgStyle} from "@angular/common";

@Component({
  selector: 'app-confirm-overlay',
  standalone: true,
  imports: [
    AddPhotoButtonComponent,
    NgIf,
    NgStyle
  ],
  templateUrl: './confirm-overlay.component.html',
  styleUrl: './confirm-overlay.component.scss'
})
export class ConfirmOverlayComponent {
  @Input() target!: string;
  @Input() photosOverlay!: WritableSignal<boolean>;

}
