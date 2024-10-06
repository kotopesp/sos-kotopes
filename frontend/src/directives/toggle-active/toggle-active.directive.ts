import {Directive, ElementRef, HostListener, Input, Renderer2} from '@angular/core';

@Directive({
  selector: '[appToggleActive]',
  standalone: true
})
export class ToggleActiveDirective {

  @Input() activeClass = 'active'; // Класс для активного элемента

  constructor(private el: ElementRef, private renderer: Renderer2) {}

  @HostListener('click')
  onClick() {
    this.deactivateOthers();
    this.activate();
  }

  private deactivateOthers() {
    const parent = this.el.nativeElement.parentNode;
    const children = parent.children;

    for (const child of children) {
      this.renderer.removeClass(child, this.activeClass);
    }
  }

  private activate() {
    this.renderer.addClass(this.el.nativeElement, this.activeClass);
  }

}
