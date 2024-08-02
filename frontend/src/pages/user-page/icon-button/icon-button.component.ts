import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';

@Component({
  selector: 'app-icon-button',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './icon-button.component.html',
  styleUrl: './icon-button.component.scss'
})
export class IconButtonComponent implements OnInit {
  @Input() label = ''
  @Input() buttonColor = ''
  @Input() textColor = ''
  @Input() icon = ''
  @Output() onClick = new EventEmitter<any>()

  iconUrl = ''
  ngOnInit() {
    this.iconUrl = `url("${this.icon}")`;
  }
  onClickButton(event: any) {
    this.onClick.emit(event)
  }
}
