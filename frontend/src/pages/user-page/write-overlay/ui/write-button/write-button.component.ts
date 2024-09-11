import {Component, Input} from '@angular/core';
import {NgStyle} from "@angular/common";


interface WriteButton {
  title: string,
  icon: string,
  buttonColor: string
}

@Component({
  selector: 'app-write-button',
  standalone: true,
  imports: [
    NgStyle
  ],
  templateUrl: './write-button.component.html',
  styleUrl: './write-button.component.scss'
})
export class WriteButtonComponent {
  @Input() WriteButton!: WriteButton;
}
