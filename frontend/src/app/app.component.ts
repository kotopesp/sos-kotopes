import {Component} from '@angular/core';
import {RouterOutlet} from '@angular/router';
import {NgForOf, NgIf} from "@angular/common";
import {StickerComponent} from "./common-ui/sticker/sticker.component";
import {QuestionComponent} from "./common-ui/question/question.component";

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [
    RouterOutlet
  ],
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
}
