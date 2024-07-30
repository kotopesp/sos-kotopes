import { Component } from '@angular/core';
import { FindStatusComponent } from "./find-status/find-status.component";
import { FindStatusBottonComponent } from "./find-status/find-status-botton/find-status-botton.component";

@Component({
  selector: 'app-filters-bar',
  standalone: true,
  imports: [FindStatusComponent, FindStatusBottonComponent],
  templateUrl: './filters-bar.component.html',
  styleUrl: './filters-bar.component.scss'
})
export class FiltersBarComponent {

}
