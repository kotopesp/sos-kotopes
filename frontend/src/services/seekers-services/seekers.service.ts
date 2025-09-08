import { Seeker, SeekerResponse } from '../../model/seeker';
import { Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { SeekerFilterService } from './seeker-filter-service/seeker-filter.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class SeekerService {
  limit: string | null = null;
  offset: string | null = null;
  private apiUrl = environment.apiUrl;

  constructor(private http: HttpClient, private filterService: SeekerFilterService) { }

  getSeekersProfile() {    
    console.log('=== ОТПРАВКА ЗАПРОСА ===');
    const filterTags = this.filterService.getTags()
    const location = this.filterService.getLocation()
    
    let params = new HttpParams()
      .set("limit", this.limit || "10")
      .set("offset", this.offset || "0");
  
    if (location) {
      params = params.set("location", location);
    }
  
    params = params.delete('animal_type');
    const catSeeker = filterTags['isCat'];
    const dogSeeker = filterTags['isDog'];
    const BothSeeker = filterTags['isBoth'];
    if ((catSeeker && dogSeeker && BothSeeker) || (!catSeeker && !dogSeeker && !BothSeeker)) {
      console.log("Показываем всех отловщиков");
    } 
    else if (BothSeeker || (catSeeker && dogSeeker)) {
      params = params.append('animal_type', 'both');
    } else if (catSeeker) {
      params = params.append('animal_type', 'cat');
    } else if (dogSeeker) {
      params = params.append('animal_type', 'dog');
    }
    if (filterTags['isMetallCage']) {
      params = params.set("have_metal_cage", true);
    }
    if (filterTags['isPlasticCage']) {
      params = params.set("have_plastic_cage", true);
    }
    if (filterTags['isNet']) {
      params = params.set("have_net", true);
    }
    if (filterTags['isLadder']) {
      params = params.set("have_ladder", true);
    }
    if (filterTags['isOther']) {
      params = params.set("have_other", filterTags['isOther']);
    }
    if (filterTags['isFree'] && !filterTags['isPay']) {
      params = params.set("max_price", 1);
    }
    if (filterTags['isPay'] && !filterTags['isFree']) {
      params = params.set("min_price", 1);
    }
    if (filterTags['isDeal']) {
      params = params.set("min_price", 0); 
      params = params.set("max_price", 0);
    }
    if (filterTags['haveCar'] && !filterTags['haventCar']) {
      params = params.set("have_car", "true");
    } else if (!filterTags['haveCar'] && filterTags['haventCar']) {
      params = params.set("have_car", "false");
    }

    console.log('Параметры запроса:', params.toString());
    return this.http.get<SeekerResponse>(`${this.apiUrl}seekers`, {params});
  }

  createSeeker(payload: FormData) {
    return this.http.post<SeekerResponse>(`${environment.apiUrl}seekers`, payload).subscribe(
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

  getSeeker(user_id: number): Observable<SeekerResponse> {
    return this.http.get<SeekerResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`)
  }
  deleteSeeker(user_id: number): Observable<SeekerResponse> {
    return this.http.delete<SeekerResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`,)
  }
  updateSeeker(updatedSeeker: Seeker): Observable<SeekerResponse> {
    const formData = new FormData();
    Object.keys(updatedSeeker).forEach(key => {
      const value = updatedSeeker[key as keyof Seeker];
      if (value !== null && value !== undefined) {
        formData.append(key, value.toString());
      }
    });
    
    return this.http.patch<SeekerResponse>(
      `${this.apiUrl}seekers/${updatedSeeker.user_id}`, 
      formData
    );
  }
}
