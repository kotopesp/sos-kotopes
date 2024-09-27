import {Component, EventEmitter, Input, Output, signal, WritableSignal} from '@angular/core';

@Component({
  selector: 'app-add-to-chat',
  standalone: true,
  imports: [],
  templateUrl: './add-to-chat.component.html',
  styleUrl: './add-to-chat.component.scss'
})
export class AddToChatComponent {
  @Input() countInArray: WritableSignal<number> = signal(0);
  isAdded = 0;
  @Input() username!: string;  
  @Input() userId!: number;

  @Output() userSelectionChanged = new EventEmitter<number>(); // для передачи изменений в родителя

  addToArray() {
    if (this.isAdded === 0) {
      this.isAdded = 1;
      this.countInArray.update(num => num + 1);
      this.userSelectionChanged.emit(this.userId);
    } else {
      this.isAdded = 0;
      this.countInArray.update(num => num - 1);
      this.userSelectionChanged.emit(-this.userId);
    }
  }
}
