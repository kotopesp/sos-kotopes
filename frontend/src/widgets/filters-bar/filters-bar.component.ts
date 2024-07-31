import { Component } from '@angular/core';
import { FindStatusBottonComponent } from "./find-status/find-status-botton/find-status-botton.component";

@Component({
  selector: 'app-filters-bar',
  standalone: true,
  imports: [FindStatusBottonComponent],
  templateUrl: './filters-bar.component.html',
  styleUrl: './filters-bar.component.scss'
})
export class FiltersBarComponent {

}
