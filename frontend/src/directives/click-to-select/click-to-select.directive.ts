import { Directive, ElementRef, HostListener, Input, Renderer2, OnInit } from '@angular/core';
import { ButtonStateService } from "../../services/button-state/button-state.service";

@Directive({
  standalone: true,
  selector: '[appClickToSelect]'
})
export class ClickToSelectDirective implements OnInit {
  @Input('appClickToSelect') classNameSelect!: string;
  @Input() buttonIndex!: number;
  private className!: string;

  constructor(
    private el: ElementRef,
    private renderer: Renderer2,
    private buttonState: ButtonStateService
  ) {}

  ngOnInit() {
    this.className = this.el.nativeElement.className;
    const savedButtonIndex = this.buttonState.getState(this.className);
    if (savedButtonIndex !== null && savedButtonIndex === this.buttonIndex) {
      this.renderer.addClass(this.el.nativeElement, this.classNameSelect);
    }
  }

  @HostListener('click') onClick() {
    const element = this.el.nativeElement;
    this.buttonState.setState(this.className, this.buttonIndex);
    const elementsWithClass = document.querySelectorAll(`.${this.classNameSelect}`);
    elementsWithClass.forEach((el) => {
      this.renderer.removeClass(el, this.classNameSelect);
    });
    if (!element.classList.contains(this.classNameSelect)) {
      this.renderer.addClass(element, this.classNameSelect);
    } else {
      this.renderer.removeClass(element, this.classNameSelect);
    }
  }
}
