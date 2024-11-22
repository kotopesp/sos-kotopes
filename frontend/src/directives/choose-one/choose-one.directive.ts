import { Directive, ElementRef, HostListener, Input, Renderer2, OnInit } from '@angular/core';
import { ButtonStateService } from "../../services/button-state/button-state.service";

@Directive({
  selector: '[appChooseOne]',
  standalone: true
})
export class ChooseOneDirective implements OnInit {
  // Статическое хранилище для экземпляров директивы
  private static activeInstances: { [key: string]: ChooseOneDirective[] } = {};
  @Input() buttonIndex!: number;

  constructor(
    private el: ElementRef,
    private renderer: Renderer2,
    private buttonState: ButtonStateService
  ) {
    const className = this.el.nativeElement.className;
    if (!ChooseOneDirective.activeInstances[className]) {
      ChooseOneDirective.activeInstances[className] = [];
    }
    ChooseOneDirective.activeInstances[className].push(this);
  }

  ngOnInit() {
    // Если в сервисе есть выбранная кнопка, обновляем состояние
    const groupKey = this.el.nativeElement.className; // Получаем класс кнопки как идентификатор группы
    const activeButtonIndex = this.buttonState.getState(groupKey);

    if (activeButtonIndex !== null) {
      // Если есть информация о нажатой кнопке, устанавливаем активное состояние
      if (this.buttonIndex === activeButtonIndex) {
        this.setActiveStyles();
      } else {
        this.setInactiveStyles();
      }
    }
  }

  @HostListener('click') onClick() {
    const className = this.el.nativeElement.className; // Получаем класс элемента

    // Сбрасываем стили для всех элементов с этим классом
    this.resetOtherElements(className);

    // Меняем стили текущего элемента
    this.setActiveStyles();

    // Сохраняем состояние выбранной кнопки в сервисе
    this.buttonState.setState(className, this.buttonIndex); // Сохраняем индекс как строку
  }

  private resetOtherElements(className: string) {
    const instances = ChooseOneDirective.activeInstances[className]; // Получаем все экземпляры по классу
    if (instances) {
      instances.forEach((instance) => {
        if (instance !== this) {
          instance.setInactiveStyles(); // Устанавливаем неактивные стили для остальных
        }
      });
    }
  }

  private setActiveStyles() {
    const link = this.el.nativeElement.querySelector('a'); // Ищем элемент <a> внутри контейнера
    if (link) {
      this.renderer.removeStyle(link, 'background-color'); // Убираем фон для активной ссылки
      // Ваши стили активного элемента
      if (link.classList.contains('animals__button__looking-for-home')) {
        this.renderer.setStyle(link, 'background-color', '#946C66'); // Цвет фона для активной ссылки
      }
      this.renderer.setStyle(link, 'border', '2px solid white'); // Белый бордер
      this.renderer.setStyle(link, 'opacity', '1'); // Устанавливаем видимость
    }
  }

  private setInactiveStyles() {
    const link = this.el.nativeElement.querySelector('a'); // Ищем элемент <a> внутри контейнера
    if (link) {
      // Ваши стили неактивного элемента
      if (link.classList.contains('animals__button__looking-for-home')) {
        this.renderer.setStyle(link, 'border', '2px solid white'); // Белый бордер
        this.renderer.removeStyle(link, 'background-color'); // Убираем фон для неактивной ссылки
      } else {
        this.renderer.setStyle(link, 'border', 'none'); // Убираем бордер
        this.renderer.setStyle(link, 'background-color', 'white'); // Белый фон для неактивной ссылки
      }
      this.renderer.setStyle(link, 'opacity', '0.3'); // Устанавливаем непрозрачность для неактивной ссылки
    }
  }
}
