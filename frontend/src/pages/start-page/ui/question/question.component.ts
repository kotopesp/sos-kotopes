import {Component, Input} from '@angular/core';
import {NgForOf, NgIf} from "@angular/common";
import {Question} from "../../../../model/question.interface";

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
  @Input() questionItem: Question
  active: string;
  notActive: string;
  hideAnswer: boolean;

  constructor() {
    this.active = "active";
    this.notActive = "";
    this.hideAnswer = false;
    this.questionItem = new class Question implements Question {
      answer: string;
      question: string;

      constructor() {
        this.answer = '';
        this.question = '';
      }
    };
  }
}
