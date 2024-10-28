import { Directive, ElementRef, HostListener, Renderer2 } from '@angular/core';

@Directive({
  selector: '[appChooseOne]',
  standalone: true
})
export class ChooseOneDirective {
  // Статическое хранилище для экземпляров директивы
  private static activeInstances: { [key: string]: ChooseOneDirective[] } = {};

  constructor(private el: ElementRef, private renderer: Renderer2) {
    // Сохраняем экземпляр директивы в статическом хранилище по классу
    const className = this.el.nativeElement.className;
    if (!ChooseOneDirective.activeInstances[className]) {
      ChooseOneDirective.activeInstances[className] = [];
    }
    ChooseOneDirective.activeInstances[className].push(this);
  }

  @HostListener('click') onClick() {
    const className = this.el.nativeElement.className; // Получаем класс элемента

    // Сбрасываем стили для всех элементов с этим классом
    this.resetOtherElements(className);

    // Меняем стили текущего элемента
    this.setActiveStyles();
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
