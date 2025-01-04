import { NgClass } from '@angular/common';
import { Trapper } from '../../../../model/trapper';
import { Component, Input } from '@angular/core';

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
    switch(this.trapper.price) {
      case -1: {
        this.divLabel = "₽ По ситуации"
        this.trapperCostStatus = "deal"
        break
      }

      case 0: {
        this.divLabel = "₽ Бесплатно"
        this.trapperCostStatus = "free"
        break
      }

      default: {
        this.divLabel = "₽ Платно"
        this.trapperCostStatus = "pay"
        break
      }
    }
  }
}
