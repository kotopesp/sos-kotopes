import {Component} from '@angular/core';
import {RouterOutlet} from "@angular/router";
import {NgForOf, NgIf} from "@angular/common";
import {StickerComponent} from "./sticker/sticker.component";
import {QuestionComponent} from "./question/question.component";
import {HeaderComponent} from "../../../widgets/header/header.component";
import {Question} from "../../../model/question.interface";

@Component({
  selector: 'app-start-page',
  standalone: true,
  imports: [RouterOutlet, NgForOf, NgIf, StickerComponent, QuestionComponent, HeaderComponent],
  templateUrl: './start-page.component.html',
  styleUrl: './start-page.component.scss'
})
export class StartPageComponent {
 
  // Временно здесь находятся тестовые
  // вопросы для проверки работоспособности блока FAQ.
  // В будущем поменяется на рабочий вариант
  // ответов на вопросы

  questionItems: Question[] = [
    {
      question: 'Рыба текст, я не знаю, что здесть писать?',
      answer: 'Рыба текст, я не знаю, что здесть писать? Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?'
    },
    {
      question: 'Вот такой вот вопрос, да, вот так?',
      answer: 'Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?'
    },
  ]
}
