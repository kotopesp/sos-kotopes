import { Injectable } from '@angular/core';
import { BehaviorSubject, Subject } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TrapperFilterService {
  public location = new BehaviorSubject("");
  public filterTags = new BehaviorSubject(this.getInitionalFilter());


  addTag(key: string, flag: boolean){
    var changedFilter = this.filterTags.getValue()
    changedFilter[key] = flag
    this.filterTags.next(changedFilter)
  }
  getTags(){
    return this.filterTags.getValue()
  }

  getInitionalFilter(){
    let flags: { [key: string]: boolean } = {};
    flags['isCat'] = false
    flags['isDog'] = false
    flags['isCadog'] = false
    flags['isMetallCage'] = false
    flags['isPlasticCage'] = false
    flags['isNet'] = false
    flags['isLadder'] = false
    flags['isOther'] = false
    flags['isPay'] = false
    flags['isFree'] = false
    flags['isDeal'] = false
    flags['haveCar'] = false
    return flags
  }

  isFilterEmpty(){
    var checker = false
    for (const key in this.filterTags) {
      checker = checker || this.filterTags.getValue()[key]
    }
    return !checker
  }

  chagneLocation(newLocation: string){
    this.location.next(newLocation)
  }
  getLocation(){
    return this.location.getValue()
  }

}
