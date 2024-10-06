import {Directive, ElementRef, HostListener, Input, Renderer2} from '@angular/core';

@Directive({
  standalone: true,
  selector: '[appClickToSelect]'
})
export class ClickToSelectDirective {
  @Input('appClickToSelect') classNameSelect!: string;

  constructor(private el: ElementRef, private renderer: Renderer2) {}

  @HostListener('click') onClick() {
    const element = this.el.nativeElement;

    // Удаляем класс со всех других элементов, у которых он есть
    const elementsWithClass = document.querySelectorAll(`.${this.classNameSelect}`);
    elementsWithClass.forEach(el => {
      this.renderer.removeClass(el, this.classNameSelect);
    });

    // Переключаем класс на текущем элементе
    if (!element.classList.contains(this.classNameSelect)) {
      this.renderer.addClass(element, this.classNameSelect);
    } else {
      this.renderer.removeClass(element, this.classNameSelect);
    }
  }
}
