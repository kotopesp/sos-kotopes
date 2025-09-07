import { Injectable } from '@angular/core';
import { BehaviorSubject, } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KeeperFilterService {
  public location = new BehaviorSubject('');
  public filterTags = new BehaviorSubject(this.getInitionalFilter());


  addTag(key: string, flag: boolean){
    const changedFilter = this.filterTags.getValue()
    changedFilter[key] = flag
    this.filterTags.next(changedFilter)
  }
  getTags(){
    return this.filterTags.getValue()
  }

  getInitionalFilter(){
    const flags: Record<string, boolean> = {}
    flags['isCat'] = false
    flags['isKitten'] = false
    flags['isDog'] = false
    flags['isPuppy'] = false
    flags['isDay'] = false
    flags['isDays'] = false
    flags['isWeeks'] = false
    flags['isMonths'] = false
    flags['isSituationallyTime'] = false
    flags['isPay'] = false
    flags['isFree'] = false
    flags['isDeal'] = false
    flags['isPerson'] = false
    flags['isOrganization'] = false
    flags['isStrays'] = false
    flags['isDomesticatedStrays'] = false
    flags['isSituationallyTaking'] = false
    flags['hasCage'] = false
    flags['hasntCage'] = false
    return flags
  }

  chagneLocation(newLocation: string){
    this.location.next(newLocation)
  }
  getLocation(){
    return this.location.getValue()
  }

}