import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root',
})
export class ButtonStateService {
  private states: { [groupKey: string]: number | null } = {};
  getState(groupKey: string): number | null {
    const value = this.states[groupKey];
    return value !== undefined ? Number(value) : null;
  }

  setState(groupKey: string, buttonIndex: number | null): void {
    this.states[groupKey] = buttonIndex;
  }
}
