import { Component } from '@angular/core';
import { NgClass } from "@angular/common";

@Component({
  selector: 'app-find-status',
  standalone: true,
  imports: [NgClass],
  templateUrl: './find-status.component.html',
  styleUrl: './find-status.component.scss'
})
export class FindStatusComponent {
  // here we get data from the backend

  status: string = "Найден"
  getClass(status: string): string {
    switch (status) {
      case 'Пропал': return 'lost_status';
      case 'Найден': return 'looking-for-home_status';
      case 'Ищет дом': return 'found-home_status';
      default: return '';
    }
  }
}
