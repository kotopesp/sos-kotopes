import {Component, EventEmitter, Output} from '@angular/core';
import {DatePipe, NgForOf, NgIf} from "@angular/common";
import {ClickToSelectDirective} from "../../../../directives/click-to-select/click-to-select.directive";

@Component({
  selector: 'app-custom-calendar',
  standalone: true,
  imports: [
    DatePipe,
    NgForOf,
    NgIf,
    ClickToSelectDirective
  ],
  templateUrl: './custom-calendar.component.html',
  styleUrl: './custom-calendar.component.scss'
})

export class CustomCalendarComponent {
  @Output() valueChange: EventEmitter<Date> = new EventEmitter<Date>();
  currentDate = new Date(); // текущая дата
  currentMonth: number;
  currentYear: number;
  daysInMonth: number[] = [];
  firstDayOffset: number = 0;
  dayNames: string[] = ['пн', 'вт', 'ср', 'чт', 'пт', 'сб', 'вс'];
  selectedDate!: Date;
  displayMonth!: string;

  constructor() {
    this.currentMonth = this.currentDate.getMonth();
    this.currentYear = this.currentDate.getFullYear();
    this.generateDaysInMonth();
  }

  generateDaysInMonth() {
    const firstDay = new Date(this.currentYear, this.currentMonth, 1).getDay(); // Определяем день недели первого дня месяца
    // В JavaScript неделя начинается с воскресенья (0), корректируем для понедельника (1)
    this.firstDayOffset = (firstDay === 0) ? 6 : firstDay - 1;
    const daysInMonth = new Date(this.currentYear, this.currentMonth + 1, 0).getDate();
    this.daysInMonth = Array.from({ length: daysInMonth }, (_, i) => i + 1);
    const displayMonth = new Date(this.currentYear, this.currentMonth);
    this.displayMonth = displayMonth.toLocaleDateString('ru-RU', { month: 'long' });
    this.displayMonth = this.displayMonth.charAt(0).toUpperCase() + this.displayMonth.slice(1)
  }

  previousMonth() {
    if (this.currentMonth === 0) {
      this.currentMonth = 11;
      this.currentYear--;
    } else {
      this.currentMonth--;
    }
    this.clearSelection();
    this.generateDaysInMonth();
  }

  nextMonth() {
    if (this.currentMonth === 11) {
      this.currentMonth = 0;
      this.currentYear++;
    } else {
      this.currentMonth++;
    }
    this.clearSelection();
    this.generateDaysInMonth();
  }

  selectDate(day: number) {
    this.selectedDate = new Date(this.currentYear, this.currentMonth, day);
  }

  // Метод для удаления класса selected у всех элементов с этим классом
  clearSelection() {
    const selectedElements = document.querySelectorAll('.selectedDay');
    selectedElements.forEach((element) => {
      element.classList.remove('selectedDay');
    });
  }

  sendData() {
    this.valueChange.emit(this.selectedDate); // Генерация события с передачей данных
  }
}
