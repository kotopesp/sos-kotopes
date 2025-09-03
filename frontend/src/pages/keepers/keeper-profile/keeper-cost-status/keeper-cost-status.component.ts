import { NgClass } from '@angular/common';
import { Component, Input } from '@angular/core';


@Component({
  selector: 'app-keeper-cost-status',
  standalone: true,
  imports: [NgClass],
  templateUrl: './keeper-cost-status.component.html',
  styleUrl: './keeper-cost-status.component.scss'
})
export class KeeperCostStatusComponent {
  @Input() keeper!: any;    // temporaly
  ngOnInit() {
    this.setFlagClass();
  }
  keeperCostStatus = ""
  divLabel = ""     
  setFlagClass() {      //temporaly
      this.divLabel = "₽ Платно"    
      this.keeperCostStatus = "pay"
  }
} 

