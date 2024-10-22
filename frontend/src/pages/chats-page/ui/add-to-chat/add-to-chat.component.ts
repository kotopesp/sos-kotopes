import {Component, Input, WritableSignal} from '@angular/core';

@Component({
  selector: 'app-add-to-chat',
  standalone: true,
  imports: [],
  templateUrl: './add-to-chat.component.html',
  styleUrl: './add-to-chat.component.scss'
})
export class AddToChatComponent {
  @Input() countInArray!: WritableSignal<number>;
  isAdded = 0

  addToArray() {
    if (this.isAdded === 0) {
      this.isAdded = 1;
      this.countInArray.update(num => num + 1);
    } else {
      this.isAdded = 0;
      this.countInArray.update(num => num - 1);
    }
  }
}
