import { Component } from '@angular/core';
import { NgClass } from '@angular/common';

@Component({
  selector: 'app-keepers-filter-bar',
  standalone: true,
  imports: [NgClass],
  templateUrl: './keepers-filter-bar.component.html',
  styleUrl: './keepers-filter-bar.component.scss'
})
export class KeepersFilterBarComponent {
  changeLocation(event: Event) {
      const inputLocation = (event.target as HTMLInputElement).value
  }
  isPressedCat = false
  pressButtonCat() {
    this.isPressedCat = !this.isPressedCat
  }
  isPressedKitten = false
  pressButtonKitten() {
    this.isPressedKitten = !this.isPressedKitten
  }
  isPressedDog = false 
  pressButtonDog() {
    this.isPressedDog = !this.isPressedDog
  }
  isPressedPuppy = false 
  pressButtonPuppy() {
    this.isPressedPuppy = !this.isPressedPuppy
  }
  isPressedDay = false 
  pressButtonDay() {
    this.isPressedDay = !this.isPressedDay
  }
  isPressedDays = false 
  pressButtonDays() {
    this.isPressedDays = !this.isPressedDays
  }
  isPressedWeeks = false 
  pressButtonWeeks() {
    this.isPressedWeeks = !this.isPressedWeeks
  }
  isPressedMonths = false 
  pressButtonMonths() {
    this.isPressedMonths = !this.isPressedMonths
  }
  isPressedOtherTime = false 
  pressButtonOtherTime() {
    this.isPressedOtherTime = !this.isPressedOtherTime
  }
  isPressedCostFree = false
  pressButtonCostFree() {
    this.isPressedCostFree = !this.isPressedCostFree
  }
  isPressedCostPay = false
  pressButtonCostPay() {
    this.isPressedCostPay = !this.isPressedCostPay
  }
  isPressedCostDeal = false
  pressButtonCostDeal() {
    this.isPressedCostDeal = !this.isPressedCostDeal
  }
  isPressedPerson= false
  pressButtonPerson() {
    this.isPressedPerson = !this.isPressedPerson
  }
  isPressedOrganization= false
  pressButtonOrganization() {
    this.isPressedOrganization = !this.isPressedOrganization
  }
  isPressedStrays = false
  pressButtonStrays() {
    this.isPressedStrays = !this.isPressedStrays
  }
  isPressedDomesticatedStrays = false
  pressButtonDomesticatedStrays() {
    this.isPressedDomesticatedStrays = !this.isPressedDomesticatedStrays
  }
  isPressedTakingSituationally= false
  pressButtonTakingSituationally() {
    this.isPressedTakingSituationally = !this.isPressedTakingSituationally
  }
  isPressedHaveCage = false
  pressButtonHaveCage() {
    this.isPressedHaveCage = !this.isPressedHaveCage  
  }
  isPressedHaventCage = false
  pressButtonHaventCage() {
    this.isPressedHaventCage = !this.isPressedHaveCage
  }
}
