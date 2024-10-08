import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output, OnInit } from '@angular/core';

@Component({
  selector: 'app-role-button',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './role-button.component.html',
  styleUrl: './role-button.component.scss'
})
export class RoleButtonComponent implements OnInit {
  @Input() label = ''
  @Input() buttonColor = ''
  @Input() icon = ''
  @Output() clicked = new EventEmitter<MouseEvent>()

  @Input() active = true;
  @Input() infoColor = '';


  anotherIconUrl = 'url("/assets/icons/arrow-down.svg")'
  iconUrl = ''
  ngOnInit() {
    this.iconUrl = `url("${this.icon}")`;
  }
  onClickButton(event: MouseEvent) {
    this.clicked.emit(event)
  }
}
