import { Component, EventEmitter, inject, Output } from '@angular/core';
import { NgClass } from "@angular/common";
import { TrapperService } from '../../../services/trapper-services/trapper.service';

@Component({
  selector: 'app-trapper-filter-bar',
  standalone: true,
  imports: [NgClass],
  templateUrl: './trapper-filter-bar.component.html',
  styleUrl: './trapper-filter-bar.component.scss'
})
export class TrapperFilterBarComponent {
  trappersService = inject(TrapperService)
  static getAllFlags: boolean[];
  isPressedCat = false
  pressButtonCat() {
    this.isPressedCat = !this.isPressedCat
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedDog= false
  pressButtonDog() {
    this.isPressedDog = !this.isPressedDog
    this.trappersService.updateTrappersData(this.getAllFlags())
  }

  isPressedCadog = false
  pressButtonCadog() {
    this.isPressedCadog = !this.isPressedCadog
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedCostFree = false
  pressButtonCostFree() {
    this.isPressedCostFree = !this.isPressedCostFree
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedCostPay = false
  pressButtonCostPay() {
    this.isPressedCostPay = !this.isPressedCostPay
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedCostDeal = false
  pressButtonCostDeal() {
    this.isPressedCostDeal = !this.isPressedCostDeal
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedMetallCatNap = false
  pressButtonMetallCatNap() {
    this.isPressedMetallCatNap = !this.isPressedMetallCatNap
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedPlasticCatNap = false
  pressButtonPlasticCatNap() {
    this.isPressedPlasticCatNap = !this.isPressedPlasticCatNap
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedNet = false
  pressButtonNet() {
    this.isPressedNet = !this.isPressedNet
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedLadder= false
  pressButtonLadder() {
    this.isPressedLadder = !this.isPressedLadder
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedOther= false
  pressButtonOther() {
    this.isPressedOther = !this.isPressedOther  
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedHaveCar= false
  pressButtonHaveCar() {
    this.isPressedHaveCar = !this.isPressedHaveCar
    this.trappersService.updateTrappersData(this.getAllFlags())
  }
  isPressedHaventCar= false
  pressButtonHaventCar() {
    this.isPressedHaventCar = !this.isPressedHaventCar
    this.trappersService.updateTrappersData(this.getAllFlags())
  }

  getAllFlags(){
    return [this.isPressedCat, this.isPressedDog, this.isPressedCadog, this.isPressedCostFree,
    this.isPressedCostPay, this.isPressedCostDeal, this.isPressedMetallCatNap, this.isPressedPlasticCatNap,
    this.isPressedNet, this.isPressedLadder, this.isPressedOther, this.isPressedHaveCar, this.isPressedHaventCar]
  }
}
