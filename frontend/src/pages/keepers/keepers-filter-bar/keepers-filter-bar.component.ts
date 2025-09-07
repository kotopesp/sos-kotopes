import { Component } from '@angular/core';
import { NgClass } from '@angular/common';
import { KeeperFilterService } from '../../../services/keepers-service/keeper-filter.service';


@Component({
  selector: 'app-keepers-filter-bar',
  standalone: true,
  imports: [NgClass],
  templateUrl: './keepers-filter-bar.component.html',
  styleUrl: './keepers-filter-bar.component.scss'
})
export class KeepersFilterBarComponent {

  constructor(private filterService: KeeperFilterService){}

  changeLocation(event: Event) {
      const inputLocation = (event.target as HTMLInputElement).value
      this.filterService.chagneLocation(inputLocation)
  }
  isPressedCat = false
  pressButtonCat() {
    this.isPressedCat = !this.isPressedCat
    this.filterService.addTag('isCat', this.isPressedCat)
  }
  isPressedKitten = false
  pressButtonKitten() {
    this.isPressedKitten = !this.isPressedKitten
    this.filterService.addTag('isKitten', this.isPressedKitten)
  }
  isPressedDog = false 
  pressButtonDog() {
    this.isPressedDog = !this.isPressedDog
    this.filterService.addTag('isDog', this.isPressedDog)
  }
  isPressedPuppy = false 
  pressButtonPuppy() {
    this.isPressedPuppy = !this.isPressedPuppy
    this.filterService.addTag('isPuppy', this.isPressedPuppy)
  }
  isPressedDay = false 
  pressButtonDay() {
    this.isPressedDay = !this.isPressedDay
    this.filterService.addTag('isDay', this.isPressedDay)
  }
  isPressedDays = false 
  pressButtonDays() {
    this.isPressedDays = !this.isPressedDays
    this.filterService.addTag('isDays', this.isPressedDays)
  }
  isPressedWeeks = false 
  pressButtonWeeks() {
    this.isPressedWeeks = !this.isPressedWeeks
    this.filterService.addTag('isWeeks', this.isPressedWeeks)
  }
  isPressedMonths = false 
  pressButtonMonths() {
    this.isPressedMonths = !this.isPressedMonths
    this.filterService.addTag('isMonths', this.isPressedMonths)
  }
  isPressedOtherTime = false 
  pressButtonOtherTime() {
    this.isPressedOtherTime = !this.isPressedOtherTime
    this.filterService.addTag('isSituationallyTime', this.isPressedOtherTime)
  }
  isPressedCostFree = false
  pressButtonCostFree() {
    this.isPressedCostFree = !this.isPressedCostFree
    this.filterService.addTag('isPay', this.isPressedCostFree)
  }
  isPressedCostPay = false
  pressButtonCostPay() {
    this.isPressedCostPay = !this.isPressedCostPay
    this.filterService.addTag('isFree', this.isPressedCostPay)
  }
  isPressedCostDeal = false
  pressButtonCostDeal() {
    this.isPressedCostDeal = !this.isPressedCostDeal
    this.filterService.addTag('isDeal', this.isPressedCostDeal)
  }
  isPressedPerson= false
  pressButtonPerson() {
    this.isPressedPerson = !this.isPressedPerson
    this.filterService.addTag('isPerson', this.isPressedPerson)
  }
  isPressedOrganization= false
  pressButtonOrganization() {
    this.isPressedOrganization = !this.isPressedOrganization
    this.filterService.addTag('isOrganization', this.isPressedOrganization)
  }
  isPressedStrays = false
  pressButtonStrays() {
    this.isPressedStrays = !this.isPressedStrays
    this.filterService.addTag('isStrays', this.isPressedStrays)
  }
  isPressedDomesticatedStrays = false
  pressButtonDomesticatedStrays() {
    this.isPressedDomesticatedStrays = !this.isPressedDomesticatedStrays
    this.filterService.addTag('isDomesticatedStrays', this.isPressedDomesticatedStrays)
  }
  isPressedTakingSituationally= false
  pressButtonTakingSituationally() {
    this.isPressedTakingSituationally = !this.isPressedTakingSituationally
    this.filterService.addTag('isSituationallyTaking', this.isPressedTakingSituationally)
  }
  isPressedHaveCage = false
  pressButtonHaveCage() {
    this.isPressedHaveCage = !this.isPressedHaveCage  
    this.filterService.addTag('isHasCage', this.isPressedHaveCage)
  }
  isPressedHaventCage = false
  pressButtonHaventCage() {
    this.isPressedHaventCage = !this.isPressedHaveCage
    this.filterService.addTag('isHasntCage', this.isPressedHaventCage)
  }
}
