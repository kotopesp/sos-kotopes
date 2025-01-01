import { Injectable } from '@angular/core';
import { environment } from '../../../environments/environment';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Trapper } from '../../../model/trapper'; 

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
          .set("isCat", filterTags[0].toString())
          .set("isDog", filterTags[1].toString())
          .set("isCadog", filterTags[2].toString())
          .set("isFree", filterTags[3].toString())
          .set("isPay", filterTags[4].toString())
          .set("isDeal", filterTags[5].toString())
          .set("isMetallCatNap", filterTags[6].toString())
          .set("isPlasticCatNap", filterTags[7].toString())
          .set("isNet", filterTags[8].toString())
          .set("isLadder", filterTags[9].toString())
          .set("isOther", filterTags[10].toString())
          .set("haveCar", filterTags[11].toString())
          .set("haventCar", filterTags[12].toString())
      }
    return this.http.get<Trapper[]>(`${this.apiUrl}/trappers`, {params})
  }
}
