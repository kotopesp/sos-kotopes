import { Component, OnInit, OnDestroy } from "@angular/core";
import { SeekerProfileComponent } from "./seeker-profile/seeker-profile.component";
import { SeekerFilterBarComponent } from "./seeker-filter-bar/seeker-filter-bar.component";
import { Meta, Seeker, SeekerResponse } from '../../model/seeker';
import { SeekerService } from '../../services/seekers-services/seekers.service';
import { SeekerFilterService } from "../../services/seekers-services/seeker-filter-service/seeker-filter.service";
import { Subject, debounceTime, switchMap, takeUntil } from 'rxjs';

@Component({
  selector: 'app-seekers',
  standalone: true,
  imports: [SeekerProfileComponent, SeekerFilterBarComponent],
  templateUrl: './seekers.component.html',
  styleUrl: './seekers.component.scss'
})
export class SeekersComponent implements OnInit, OnDestroy {
  seekers: Seeker[] = [];
  meta: Meta | null = null;

  private destroy$ = new Subject<void>();
  private filterChange$ = new Subject<void>();

  constructor(
    private seekerService: SeekerService,
    private filterService: SeekerFilterService
  ) {
    this.filterChange$
      .pipe(
        debounceTime(300),
        switchMap(() => this.seekerService.getSeekersProfile()), 
        takeUntil(this.destroy$) 
      )
      .subscribe({
        next: (response) => this.handleResponse(response),
        error: (error) => console.error('Ошибка запроса:', error)
      });
  }

  updateSeekersArray() {
    this.filterChange$.next();
  }


  private handleResponse(response: SeekerResponse) {
    if (response && response.data) {
      this.seekers = response.data.payload;
      this.meta = response.data.meta;
      console.log('Получено отловщиков:', this.seekers.length);
    } else {
      console.log('Нет данных в ответе');
      this.seekers = [];
    }
  }

  ngOnInit(): void {
    this.filterService.filterTags
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => this.updateSeekersArray());

    this.filterService.location
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => this.updateSeekersArray());
  }


  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}