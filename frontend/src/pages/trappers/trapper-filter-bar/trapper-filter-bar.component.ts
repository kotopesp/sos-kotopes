import { Component, } from '@angular/core';
import { NgClass } from "@angular/common";
import { TrapperFilterService } from '../../../services/trapper-services/trapper-filter-service/trapper-filter.service';


@Component({
  selector: 'app-trapper-filter-bar',
  standalone: true,
  imports: [NgClass],
  templateUrl: './trapper-filter-bar.component.html',
  styleUrl: './trapper-filter-bar.component.scss'
})
export class TrapperFilterBarComponent {

  constructor(private filterService: TrapperFilterService){}

  isPressedCat = false
  changeLocation(event: Event) {
      let inputLocation = (event.target as HTMLInputElement).value
      this.filterService.chagneLocation(inputLocation)
  }
  pressButtonCat() {
    this.isPressedCat = !this.isPressedCat
    this.filterService.addTag('isCat', this.isPressedCat)
  }
  isPressedDog= false 
  pressButtonDog() {
    this.isPressedDog = !this.isPressedDog
    this.filterService.addTag('isDog', this.isPressedDog)
  }

  isPressedCadog = false
  pressButtonCadog() {
    this.isPressedCadog = !this.isPressedCadog
    this.filterService.addTag('isCadog', this.isPressedCadog)
  }
  isPressedCostFree = false
  pressButtonCostFree() {
    this.isPressedCostFree = !this.isPressedCostFree
    this.filterService.addTag('isFree', this.isPressedCostFree)
  }
  isPressedCostPay = false
  pressButtonCostPay() {
    this.isPressedCostPay = !this.isPressedCostPay
    this.filterService.addTag('isPay', this.isPressedCostPay)
  }
  isPressedCostDeal = false
  pressButtonCostDeal() {
    this.isPressedCostDeal = !this.isPressedCostDeal
    this.filterService.addTag('isDeal', this.isPressedCostDeal)
  }
  isPressedMetallCatNap = false
  pressButtonMetallCatNap() {
    this.isPressedMetallCatNap = !this.isPressedMetallCatNap
    this.filterService.addTag('isMetallNap', this.isPressedMetallCatNap)
  }
  isPressedPlasticCatNap = false
  pressButtonPlasticCatNap() {
    this.isPressedPlasticCatNap = !this.isPressedPlasticCatNap
    this.filterService.addTag('isPlasticNap', this.isPressedPlasticCatNap)
  }
  isPressedNet = false
  pressButtonNet() {
    this.isPressedNet = !this.isPressedNet
    this.filterService.addTag('isNet', this.isPressedNet)
  }
  isPressedLadder= false
  pressButtonLadder() {
    this.isPressedLadder = !this.isPressedLadder
    this.filterService.addTag('isLadder', this.isPressedLadder)
  }
  isPressedOther= false
  pressButtonOther() {
    this.isPressedOther = !this.isPressedOther  
    this.filterService.addTag('isOther', this.isPressedOther)
  }
  isPressedHaveCar = false
  pressButtonHaveCar() {
    this.isPressedHaveCar = !this.isPressedHaveCar  
    this.filterService.addTag('haveCar', this.isPressedHaveCar)
  }
  isPressedHaventCar = false
  pressButtonHaventCar() {
    this.isPressedHaventCar = !this.isPressedHaveCar  
    this.filterService.addTag('haveCar', this.isPressedHaventCar)
  }

  /**/
}
