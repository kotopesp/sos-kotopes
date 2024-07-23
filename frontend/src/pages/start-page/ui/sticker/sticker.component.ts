import {Component, Input} from '@angular/core';
import {NgForOf, NgOptimizedImage} from "@angular/common";

@Component({
  selector: 'app-sticker',
  standalone: true,
  imports: [
    NgForOf,
    NgOptimizedImage
  ],
  templateUrl: './sticker.component.html',
  styleUrl: './sticker.component.scss'
})
export class StickerComponent {
  @Input() stickerItem: any
}
