import {Component, Input} from '@angular/core';
import {NgForOf, NgOptimizedImage} from "@angular/common";
import {Sticker} from "../../../../model/sticker.interface";

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
  @Input() stickerItem: Sticker

  constructor() {
    this.stickerItem = new class Sticker implements Sticker {
      class: string;
      icon: string;
      subtitle1: string;
      subtitle2: string;
      title: string;

      constructor() {
        this.class = '';
        this.icon = '';
        this.subtitle1 = '';
        this.subtitle2 = '';
        this.title = '';
      }
    }
  }
}
