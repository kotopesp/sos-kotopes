import { NgClass } from '@angular/common';
import { Component, Input } from '@angular/core';
import { Trapper } from '../../../../model/trapper';

@Component({
  selector: 'app-trapper-cost-status',
  standalone: true,
  imports: [NgClass],
  templateUrl: './trapper-cost-status.component.html',
  styleUrl: './trapper-cost-status.component.scss'
})
export class TrapperCostStatusComponent {
  @Input() trapper!: Trapper;
  ngOnInit() {
    this.setFlagClass();
  }

  trapperCostStatus = ""
  divLabel = ""
  setFlagClass() {
    switch(this.trapper.cost) {
      case "deal": {
        this.divLabel = "₽ По ситуации"
        this.trapperCostStatus = "deal"
        break
      }

      case "free": {
        this.divLabel = "₽ Бесплатно"
        this.trapperCostStatus = "free"
        break
      }

      case "pay": {
        this.divLabel = "₽ Платно"
        this.trapperCostStatus = "pay"
        break
      }
    }
  }
}
