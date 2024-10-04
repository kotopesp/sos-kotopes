import {Component} from '@angular/core';
import {RouterOutlet} from "@angular/router";
import {NgForOf, NgIf} from "@angular/common";
import {StickerComponent} from "./sticker/sticker.component";
import {QuestionComponent} from "./question/question.component";
import {HeaderComponent} from "../../../widgets/header/header.component";
import {Question} from "../../../model/question.interface";
import {Sticker} from "../../../model/sticker.interface";
import {ButtonFindPetComponent} from "../../../shared/buttons/button-find-pet/button-find-pet.component";
import {ButtonLostPetComponent} from "../../../shared/buttons/button-lost-pet/button-lost-pet.component";
import {
  ButtonLookingForHomeComponent
} from "../../../shared/buttons/button-looking-for-home/button-looking-for-home.component";

@Component({
  selector: 'app-start-page',
  standalone: true,
  imports: [RouterOutlet, NgForOf, NgIf, StickerComponent, QuestionComponent, HeaderComponent, ButtonFindPetComponent, ButtonLostPetComponent, ButtonLookingForHomeComponent],
  templateUrl: './start-page.component.html',
  styleUrl: './start-page.component.scss'
})
export class StartPageComponent {
  stickerItems: Sticker[] = [
    {
      class: 'yellow-color',
      title: 'Потеряшка',
      subtitle1: 'Нашли бездомное или потеряли домашнее животное?',
      subtitle2: 'Здесь вы найдете подробные инструкции и контакты ' +
        'для связи с отловщиками, частными ветеринарами и службами ' +
        'передержки. Мы поможем вам обеспечить безопасность и защиту ' +
        'для каждого найденного или пропавшего питомца',
      icon: 'paw-icon.svg'
    },
    {
      class: 'pink-color',
      title: 'Передержка',
      subtitle1: 'Нужен ночлег или помощь?',
      subtitle2: 'Если вы готовы рассказать о нас,' +
        ' предоставить ночлег для животного' +
        ' или вы частный ветеринар, или ' +
        'отловщик, то смело нажимайте сюда и ' +
        'присоединяйтесь к нашей команде!',
      icon: 'hand-icon.svg'
    },
    {
      class: 'purple-color',
      title: 'Отловщики',
      subtitle1: 'Нужен отлов улов улёт? жду когда Ангелина придумает текст',
      subtitle2: 'Если вы готовы рассказать о нас, ' +
        'предоставить ночлег для животного или вы ' +
        'частный ветеринар, или отловщик, то смело ' +
        'нажимайте сюда и присоединяйтесь к нашей команде!',
      icon: 'net-icon.svg'
    },
    {
      class: 'green-color',
      title: 'Ветеринарные Клиники',
      subtitle1: 'Ищете ближайшую ветеринарную клинику?',
      subtitle2: 'Здесь вы найдете список ближайших ' +
        'ветеринарных клиник с указанием стоимости ' +
        'услуг и льготных условий. Мы стремимся ' +
        'сделать ветеринарную помощь доступной для ' +
        'каждого животного. Быстрая и качественная ' +
        'помощь – залог здоровья наших питомцев',
      icon: 'cross-icon.svg'
    },
    {
      class: 'blue-color',
      title: 'Правовые Вопросы',
      subtitle1: 'Не знаете, что можно и чего нельзя делать?',
      subtitle2: 'Мы поможем вам разобраться в ' +
        'правовых аспектах и взаимодействии ' +
        'с органами власти. Узнайте, как ' +
        'действовать в рамках закона, что ' +
        'делать в случае жестокого обращения с ' +
        'животными и как взаимодействовать с ЖКХ ' +
        'и районными администрациями. Получите ' +
        'ответы на все важные вопросы и защитите ' +
        'права животных',
      icon: 'lawyer-icon.svg'
    }
  ]


  // Временно здесь находятся тестовые
  // вопросы для проверки работоспособности блока FAQ.
  // В будущем поменяется на рабочий вариант
  // ответов на вопросы

  questionItems: Question[] = [
    {
      question: 'Рыба текст, я не знаю, что здесть писать?',
      answer: 'Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?'
    },
    {
      question: 'Вот такой вот вопрос, да, вот так?',
      answer: 'Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?Рыба текст, я не знаю, что здесть писать?'
    },
  ]
}
