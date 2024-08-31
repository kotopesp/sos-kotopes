import {Component, Input} from '@angular/core';
import {NgStyle} from "@angular/common";

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
  @Input() WriteButton: any;
}
