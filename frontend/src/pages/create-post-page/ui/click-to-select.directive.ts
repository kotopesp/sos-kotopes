import {Directive, ElementRef, HostListener, Input, Renderer2} from '@angular/core';

@Directive({
  standalone: true,
  selector: '[appClickToSelect]'
})
export class ClickToSelectDirective {
  // Статическая переменная для хранения ранее выбранного элемента
  private static selectedElement: HTMLElement | null = null;

  constructor(private el: ElementRef) {}

  @HostListener('click') onClick() {
    const currentElement = this.el.nativeElement;

    // Если уже есть выбранный элемент и он отличается от текущего
    if (ClickToSelectDirective.selectedElement && ClickToSelectDirective.selectedElement !== currentElement) {
      // Снимаем класс выделения с предыдущего элемента
      ClickToSelectDirective.selectedElement.classList.remove('selected');
    }

    // Устанавливаем или убираем класс 'selected' для текущего элемента
    currentElement.classList.toggle('selected');

    // Обновляем статическую переменную для отслеживания выбранного элемента
    ClickToSelectDirective.selectedElement = currentElement.classList.contains('selected')
      ? currentElement
      : null;
  }
}
