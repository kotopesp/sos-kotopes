import { Injectable } from '@angular/core';
import { BehaviorSubject, } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SeekerFilterService {
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
    flags['isDog'] = false
    flags['isBoth'] = false
    flags['isMetallCage'] = false
    flags['isPlasticCage'] = false
    flags['isNet'] = false
    flags['isLadder'] = false
    flags['isOther'] = false
    flags['isPay'] = false
    flags['isFree'] = false
    flags['isDeal'] = false
    flags['haveCar'] = false
    flags['haventCar'] = false
    return flags
  }

  chagneLocation(newLocation: string){
    this.location.next(newLocation)
  }
  getLocation(){
    return this.location.getValue()
  }

}
