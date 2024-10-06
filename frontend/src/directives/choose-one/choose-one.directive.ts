import {Directive, ElementRef, HostListener, Input, Renderer2} from '@angular/core';

@Directive({
  selector: '[appChooseOne]',
  standalone: true
})
export class ChooseOneDirective {
  @Input('appChooseOne') activeClass!: string
  inactiveClass: string;

  private static activeElement: ChooseOneDirective | null = null; // Хранит активный элемент

  constructor(private el: ElementRef, private renderer: Renderer2) {
    this.inactiveClass = '';
  }



  @HostListener('click') onClick() {
    // Если есть активный элемент, сбрасываем его стили
    if (ChooseOneDirective.activeElement) {
      ChooseOneDirective.activeElement.setInactiveStyles();
    }

    // Устанавливаем текущий элемент как активный
    ChooseOneDirective.activeElement = this;

    // Меняем стили текущего элемента
    this.setActiveStyles();
  }

  private setActiveStyles() {
    this.renderer.setStyle(this.el.nativeElement, 'border', '2px solid white'); // Стиль для активного элемента
  }

  private setInactiveStyles() {
    this.renderer.setStyle(this.el.nativeElement, 'background-color', 'white'); // Стиль для неактивных элементов
    this.renderer.setStyle(this.el.nativeElement, 'color', '#443416'); // Другие свойства для неактивных элементов
    this.renderer.setStyle(this.el.nativeElement, 'opacity', '0.3'); // Другие свойства для неактивных элементов
  }
}
