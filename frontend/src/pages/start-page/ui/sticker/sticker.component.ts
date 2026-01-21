import {Component, Input} from '@angular/core';
import {NgForOf, NgOptimizedImage} from "@angular/common";

@Component({
  selector: 'app-sticker',
  standalone: true,
  imports: [NgForOf, NgOptimizedImage],
  templateUrl: './sticker.component.html',
  styleUrl: './sticker.component.scss',
})
export class StickerComponent {
  @Input() classEl = '';
  @Input() classColor = '';
  @Input() icon = '';
  @Input() subtitle1 = '';
  @Input() subtitle2 = '';
  @Input() title = '';
}