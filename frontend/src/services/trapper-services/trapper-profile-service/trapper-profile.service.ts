;
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../../environments/environment';
import { Trapper } from '../../../model/trapper'; 
import { Injectable } from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class TrapperProfileService {
  
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient) { }

  getTrappersProfile(filterTags: boolean[]) {
    const params = new HttpParams()
    if (filterTags.length != 0){
      const params = new HttpParams()
          .set("animal_type", filterTags[0].toString())
          .set("location ", filterTags[1].toString())
          .set("min_equipment_rental", filterTags[2].toString())
          .set("max_equipment_rental", filterTags[3].toString())
          .set("have_metal_cage", filterTags[4].toString())
          .set("have_plastic_cage", filterTags[5].toString())
          .set("have_net", filterTags[6].toString())
          .set("have_ladder", filterTags[7].toString())
          .set("have_other", filterTags[8].toString())
          .set("have_car", filterTags[9].toString())
          .set("isOther", filterTags[10].toString())
          .set("haveCar", filterTags[11].toString())
          .set("haventCar", filterTags[12].toString())
      }
    return this.http.get<Trapper[]>(`${this.apiUrl}/seekers`, {params})
  }
}
