
import { TrapperProfileService } from './trapper-profile-service/trapper-profile.service';
import { Trapper } from '../../model/trapper';
import { inject, Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class TrapperService {
  trapperProfileService = inject(TrapperProfileService)
  trappers: Trapper[] = []
  constructor() {this.updateTrappersData([])}
  updateTrappersData(filter: boolean[]){
    this.trapperProfileService.getTrappersProfile(filter).subscribe(val => {
      this.trappers = val
    })
  }
}
