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
  @Input() selectedFiles!: { name: string, preview: string }[];
  @Input() thisOverlay!: WritableSignal<boolean>;
  @Input() numberOfSlide!: WritableSignal<number>;

  goToNext() {
    this.thisOverlay.set(false);
    this.numberOfSlide.set(this.numberOfSlide() + 1)
  }
}
