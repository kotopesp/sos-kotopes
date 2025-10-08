import { Keeper, KeeperResponse } from '../../model/keeper';
import { Inject, Injectable } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { environment } from '../../environments/environment';
import { KeeperFilterService } from './keeper-filter.service';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class KeeperService {
  limit: string | null = null;
  offset: string | null = null;
  private apiUrl = environment.apiUrl;

  constructor(
    private http: HttpClient, 
    @Inject(KeeperFilterService) private filterService: KeeperFilterService
  ) { }


  getKeepersProfile() {    
    console.log('=== ОТПРАВКА ЗАПРОСА ===');
    const filterTags = this.filterService.getTags()
    const location = this.filterService.getLocation()
    
    let params = new HttpParams()
      .set("limit", this.limit || "10")
      .set("offset", this.offset || "0");
  
    if (location) {
      params = params.set("location", location);
    }
  
    params = params.delete('animal_category');
    const catKeeper = filterTags['isCat'];
    const dogKeeper = filterTags['isDog'];

    if ((catKeeper && dogKeeper) || (!catKeeper && !dogKeeper)) {
      console.log("Показываем всех отловщиков");
    } 
    else if (catKeeper) {
      params = params.append('animal_category', 'cat');
    } else if (dogKeeper) {
      params = params.append('animal_category', 'dog');
    }

    if (filterTags['isHours']) {
      params = params.set("boarding_duration", 'hours');
    }
    if (filterTags['isDays']) {
      params = params.set("boarding_duration", 'days');
    }
    if (filterTags['isWeeks']) {
      params = params.set("boarding_duration", 'weeks');
      console.log("Chose WEEKS!!!!")
    }
    if (filterTags['isMonths']) {
      params = params.set("boarding_duration", 'months');
    }
    if (filterTags['isSituationallyTime']) {
      params = params.set("boarding_duration", 'depends');
    }
    
    if (filterTags['isFree'] && !filterTags['isPay']) {
      params = params.set("max_price", 0);
      
    }
    if (filterTags['isPay'] && !filterTags['isFree']) {
      params = params.set("min_price", 1);
    }
    if (filterTags['isDeal']) {
      params = params.set("min_price", 0); 
      params = params.set("max_price", 0);
    }
    if (filterTags['isStrays']) {
      params = params.set("animal_acceptance", 'homeless');
      console.log("chose strays!", params.toString());
    }
    if (filterTags['isDomesticatedStrays']) {
      params = params.set("animal_acceptance", 'homeless-hadhome');
      console.log("chose homeless-hadhome!", params.toString());
    }
    if (filterTags['isSituationallyTaking']) {
      params = params.set("animal_acceptance", 'depends');
      console.log("chose depends!", params.toString());
    }
    if (filterTags['hasCage'] && !filterTags['hasntCage']) {
      console.log("chose CAGE!", params.toString());
      params = params.set("has_cage", "true");
    } else if (!filterTags['hasCage'] && filterTags['hasntCage']) {
      params = params.set("has_cage", "false");
    }

    console.log('Параметры запроса:', params.toString());
    return this.http.get<KeeperResponse>(`${this.apiUrl}keepers`, {params});
  }

  createKeeper(payload: FormData) {
    return this.http.post<KeeperResponse>(`${environment.apiUrl}keepers`, payload).subscribe(
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

  getKeeper(user_id: number): Observable<KeeperResponse> {
    return this.http.get<KeeperResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`)
  }
  deleteKeeper(user_id: number): Observable<KeeperResponse> {
    return this.http.delete<KeeperResponse>(`${environment.apiUrl}seekers/${user_id.toString()}`,)
  }
  updateKeeper(updatedSeeker: Keeper): Observable<KeeperResponse> {
    const formData = new FormData();
    Object.keys(updatedSeeker).forEach(key => {
      const value = updatedSeeker[key as keyof Keeper];
      if (value !== null && value !== undefined) {
        formData.append(key, value.toString());
      }
    });
    
    return this.http.patch<KeeperResponse>(
      `${this.apiUrl}seekers/${updatedSeeker.user_id}`, 
      formData
    );
  }
}