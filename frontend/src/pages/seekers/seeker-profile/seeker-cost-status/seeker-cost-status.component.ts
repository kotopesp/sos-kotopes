import { NgClass } from '@angular/common';
import { Seeker } from '../../../../model/seeker';
import { Component, Input, OnChanges, SimpleChanges } from '@angular/core';

@Component({
  selector: 'app-seeker-cost-status',
  standalone: true,
  imports: [NgClass],
  templateUrl: './seeker-cost-status.component.html',
  styleUrl: './seeker-cost-status.component.scss'
})
export class SeekerCostStatusComponent implements OnChanges {
  @Input() seeker!: Seeker;
  
  seekerCostStatus = "";
  divLabel = "";

  ngOnChanges(changes: SimpleChanges) {
    if (changes['seeker'] && this.seeker) {
      this.setFlagClass();
    }
  }

  setFlagClass() {
    if (!this.seeker) {
      this.divLabel = "₽ Бесплатно";
      this.seekerCostStatus = "free";
      return;
    }
    switch(this.seeker.price) {
      case -1: {
        this.divLabel = "₽ По ситуации";
        this.seekerCostStatus = "deal";
        break;
      }
      case 0: {
        this.divLabel = "₽ Бесплатно";
        this.seekerCostStatus = "free";
        break;
      }
      default: {
        this.divLabel = `₽ ${this.seeker.price}`; 
        this.seekerCostStatus = "pay";
        break;
      }
    }
  }
}