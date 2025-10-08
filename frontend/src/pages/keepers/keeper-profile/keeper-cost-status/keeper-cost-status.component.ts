import { NgClass } from '@angular/common';
import { Component, Input, SimpleChanges } from '@angular/core';
import { Keeper } from '../../../../model/keeper';


@Component({
  selector: 'app-keeper-cost-status',
  standalone: true,
  imports: [NgClass],
  templateUrl: './keeper-cost-status.component.html',
  styleUrl: './keeper-cost-status.component.scss'
})
export class KeeperCostStatusComponent {
  @Input() keeper!: Keeper;
  
  keeperCostStatus = "";
  divLabel = "";

  ngOnChanges(changes: SimpleChanges) {
    if (changes['keeper'] && this.keeper) {
      this.setFlagClass();
    }
  }

  setFlagClass() {
    if (!this.keeper) {
      this.divLabel = "₽ Бесплатно";
      this.keeperCostStatus = "free";
      return;
    }
    switch(this.keeper.price) {
      case -1: {
        this.divLabel = "₽ По ситуации";
        this.keeperCostStatus = "deal";
        break;
      }
      case 0: {
        this.divLabel = "₽ Бесплатно";
        this.keeperCostStatus = "free";
        break;
      }
      default: {
        this.divLabel = `₽ ${this.keeper.price}`; 
        this.keeperCostStatus = "pay";
        break;
      }
    }
  }
} 

