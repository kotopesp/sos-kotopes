import { Component, } from '@angular/core';
import { NgClass } from "@angular/common";
import { SeekerFilterService } from '../../../services/seekers-services/seeker-filter-service/seeker-filter.service';


@Component({
  selector: 'app-seeker-filter-bar',
  standalone: true,
  imports: [NgClass],
  templateUrl: './seeker-filter-bar.component.html',
  styleUrl: './seeker-filter-bar.component.scss'
})
export class SeekerFilterBarComponent {

  constructor(private filterService: SeekerFilterService){}

  isPressedCat = false
  changeLocation(event: Event) {
      const inputLocation = (event.target as HTMLInputElement).value
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

  isPressedBoth = false
  pressButtonBoth() {
    this.isPressedBoth = !this.isPressedBoth
    this.filterService.addTag('isBoth', this.isPressedBoth)
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
    this.filterService.addTag('isMetallCage', this.isPressedMetallCatNap)
  }
  isPressedPlasticCatNap = false
  pressButtonPlasticCatNap() {
    this.isPressedPlasticCatNap = !this.isPressedPlasticCatNap
    this.filterService.addTag('isPlasticCage', this.isPressedPlasticCatNap)
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
    this.isPressedHaventCar = !this.isPressedHaventCar  
    this.filterService.addTag('haventCar', this.isPressedHaventCar)
  }
}
