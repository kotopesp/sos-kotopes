import { Component } from '@angular/core';
import { FindStatusComponent } from "./find-status/find-status.component";

@Component({
  selector: 'app-filters-bar',
  standalone: true,
  imports: [FindStatusComponent],
  templateUrl: './filters-bar.component.html',
  styleUrl: './filters-bar.component.scss'
})
export class FiltersBarComponent {

}
