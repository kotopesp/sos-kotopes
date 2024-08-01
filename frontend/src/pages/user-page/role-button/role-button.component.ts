import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-role-button',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './role-button.component.html',
  styleUrl: './role-button.component.scss'
})
export class RoleButtonComponent {
  @Input() label = ''
  @Input() buttonColor = ''
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
