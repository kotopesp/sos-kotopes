import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ButtonStateService {
  private states: { [groupKey: string]: number | null } = {};

  // Получить состояние кнопки по группе
  getState(groupKey: string): number | null {
    // Преобразуем строку в число, если она есть, иначе возвращаем null
    const value = this.states[groupKey];
    return value !== undefined ? Number(value) : null;
  }

  setState(groupKey: string, buttonIndex: number): void {
    this.states[groupKey] = buttonIndex;  // Сохраняем как число
  }
}
