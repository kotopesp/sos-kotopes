import { Component, OnInit, OnDestroy } from "@angular/core";
import { KeeperProfileComponent } from "./keeper-profile/keeper-profile.component";
import { KeepersFilterBarComponent } from "./keepers-filter-bar/keepers-filter-bar.component";
import { Meta, Keeper, KeeperResponse } from '../../model/keeper';
import { KeeperService } from '../../services/keepers-service/keeper.service';
import { KeeperFilterService } from "../../services/keepers-service//keeper-filter.service";
import { Subject, debounceTime, switchMap, takeUntil } from 'rxjs';
import { NgIf } from "@angular/common";

@Component({
  selector: 'app-keepers',
  standalone: true,
  imports: [KeeperProfileComponent , KeepersFilterBarComponent, NgIf],
  templateUrl: './keepers.component.html',
  styleUrl: './keepers.component.scss'
})
export class KeepersComponent implements OnInit, OnDestroy  {
  keepers: Keeper[] = [];
  meta: Meta | null = null;

  private destroy$ = new Subject<void>();
  private filterChange$ = new Subject<void>();

  constructor(
    private keeperService: KeeperService,
    private filterService: KeeperFilterService
  ) {
    this.filterChange$
      .pipe(
        debounceTime(300),
        switchMap(() => this.keeperService.getKeepersProfile()), 
        takeUntil(this.destroy$) 
      )
      .subscribe({
        next: (response) => this.handleResponse(response),
        error: (error) => console.error('Ошибка запроса:', error)
      });
  }

  updateKeepersArray() {
    this.filterChange$.next();
  }


  private handleResponse(response: KeeperResponse) {
    if (response && response.data) {
      this.keepers = response.data.payload;
      this.meta = response.data.meta;
      console.log('Получено передержщиков:', this.keepers.length);
    } else {
      console.log('Нет данных в ответе');
      this.keepers = [];
    }
  }

  ngOnInit(): void {
    this.filterService.filterTags
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => this.updateKeepersArray());

    this.filterService.location
      .pipe(takeUntil(this.destroy$))
      .subscribe(() => this.updateKeepersArray());
  }


  ngOnDestroy(): void {
    this.destroy$.next();
    this.destroy$.complete();
  }
}