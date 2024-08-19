import { Component } from '@angular/core';
import { FindStatusButtonsComponent } from "./find-status/find-status-buttons/find-status-buttons.component";
import {NgForOf} from "@angular/common";
import {ColorFilter} from "../../model/color-filter.interface";

@Component({
  selector: 'app-filters-bar',
  standalone: true,
  imports: [FindStatusButtonsComponent, NgForOf],
  templateUrl: './filters-bar.component.html',
  styleUrl: './filters-bar.component.scss'
})
export class FiltersBarComponent {

  colorFilterItems: ColorFilter[] = [
    {
      class: 'black_button',
      title: 'Чёрный'
    },
    {
      class: 'white_button',
      title: 'Белый'
    },
    {
      class: 'brown_button',
      title: 'Коричневый'
    },
    {
      class: 'russet_button',
      title: 'Рыжий'
    },
    {
      class: 'gray_button',
      title: 'Серый',
    },
    {
      class: 'creamy_button',
      title: 'Кремовый'
    }
  ]
}
