import {Component} from "@angular/core";
import {ChatTypeButtonComponent} from "../../entities/chat-type-button/chat-type-button.component";

interface Button {
  label: string,
  color: string,
  iconName: string,
  counter?: number,
}

@Component({
  selector: "test-page",
  imports: [
    ChatTypeButtonComponent]
  ,
  standalone: true,
  template: `
    @for (button of buttons; track button.label) {
      <chat-type-button [label]="button.label" [color]="button.color" [counter]="button.counter || 0" [icon]="button.iconName"></chat-type-button>
    }
  `
})
export class TestPageComponent {
  buttons: Button[] = [
    {label: "Все чаты", color: "#352B1A", iconName: "allChats.svg"},
    {label: "Отклики", color: "#352303", counter: 5, iconName: "respond.svg"},
    {label: "Чаты с передержкой", color: "#2B1800", counter: 4, iconName: "keepers.svg"},
    {label: "Чаты с отловщиками", color: "#221630", iconName: "seekers.svg"},
    {label: "Чаты с ветеринарами", color: "#182C2A", iconName: "vets.svg"},
    {label: "Другие чаты", color: "#21190B", iconName: "other.svg"},
  ]
}
