import { Trapper, TrapperResponse } from '../../model/trapper';
import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { TrapperFilterService } from './trapper-filter-service/trapper-filter.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class TrapperService {
  /*trapperProfileService = inject(TrapperProfileService)
  trappers: Trapper[] = []*/

  limit: string | null = null;
  offset: string | null = null;
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient, private filterService: TrapperFilterService) { }

  getTrappersProfile() {    
    const filterTags = this.filterService.getTags()
    const location = this.filterService.getLocation()
    const params = new HttpParams()
    .set("limit", this.limit || "10")
    .set("offset", this.offset || "0")
    if (this.filterService.isFilterEmpty()){
      params.set("animal_type", this.getAnimalType())
          .set("location ", location)
          .set("have_metal_cage", filterTags['isMetallCage'])
          .set("have_plastic_cage", filterTags['isPlasticCage'])
          .set("have_net", filterTags['isNet'])
          .set("have_ladder", filterTags["isLadder"])
          .set("have_other", filterTags["isOther"])
          .set("max_price", Infinity) // TODO (rewrite backend)
          .set("min_price", 0)        // TODO (rewrite backend)
          .set("haveCar", filterTags["haveCar"])
          return this.http.get<TrapperResponse>(`${this.apiUrl}seekers`, {params})
      }
    
    return this.http.get<TrapperResponse>(`${this.apiUrl}seekers`, {params}) 
  }

  createTrapper(payload: FormData) {
    return this.http.post<TrapperResponse>(`${environment.apiUrl}seekers`, payload).subscribe(
      {
        next: () => {
          console.log('success') // TO DO: add routing to form with registration
        },
        error: (error) => {
          console.log(error);
        }
      }
    )
  }

  getTrapper(user_id: number): Observable<TrapperResponse> {
    return this.http.get<TrapperResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`)
  }
  deleteTrapper(user_id: number): Observable<TrapperResponse> {
    return this.http.delete<TrapperResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`,)
  }
  updateTrapper(updatedTrapper: Trapper): Observable<TrapperResponse> {
    const params = new HttpParams()
    params.set("animal_type", updatedTrapper.animal_type)
      .set("description", updatedTrapper.description)
      .set("location ", updatedTrapper.location)
      .set("equipment_rental", updatedTrapper.equipment_rental)
      .set("have_metal_cage",updatedTrapper.have_metal_cage)
      .set("have_plastic_cage", updatedTrapper.have_plastic_cage)
      .set("have_net", updatedTrapper.have_net)
      .set("have_ladder", updatedTrapper.have_ladder)
      .set("have_other", updatedTrapper.have_other)
      .set("haveCar", updatedTrapper.have_car)
      .set("price", updatedTrapper.price)
      .set("haveCar", updatedTrapper.have_car)
      .set("willingness_carry", updatedTrapper.willingness_carry)
    return this.http.patch<TrapperResponse>(`${this.apiUrl}seekers/${updatedTrapper.user_id.toString()}`, {params})
  }

  getAnimalType(){
    if (this.filterService.getTags()["isCat"] == true){
      return 'cat'
    }
    return "dog"
  }



}
