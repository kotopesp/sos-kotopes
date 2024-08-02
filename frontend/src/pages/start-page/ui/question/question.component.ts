import {Component, Input} from '@angular/core';
import {NgForOf, NgIf} from "@angular/common";

@Component({
  selector: 'app-question',
  standalone: true,
  imports: [
    NgIf,
    NgForOf
  ],
  templateUrl: './question.component.html',
  styleUrl: './question.component.scss'
})
export class QuestionComponent {
  @Input() questionItem: any
  active = "active";
  notActive = "";
  hideAnswer = false;
}
