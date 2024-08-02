import { Component } from '@angular/core';
import { FindStatusButtonsComponent } from "./find-status/find-status-buttons/find-status-buttons.component";

@Component({
  selector: 'app-filters-bar',
  standalone: true,
  imports: [FindStatusButtonsComponent],
  templateUrl: './filters-bar.component.html',
  styleUrl: './filters-bar.component.scss'
})
export class FiltersBarComponent {

}
