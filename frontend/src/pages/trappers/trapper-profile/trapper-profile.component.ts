import { FavoritesButtonComponent } from '../../../shared/post/favorites-button/favorites-button.component';
import { FindStatusFlagComponent } from "../../../shared/post/find-status-flag/find-status-flag.component";
import { Trapper } from '../../../model/trapper';
import { TrapperTypePipe } from "../../../pipes/trapper-type.pipe";
import { TrapperCostStatusComponent } from "./trapper-cost-status/trapper-cost-status.component";
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-trapper-profile',
  standalone: true,
  imports: [FavoritesButtonComponent, FindStatusFlagComponent, TrapperTypePipe, TrapperCostStatusComponent],
  templateUrl: './trapper-profile.component.html',
  styleUrl: './trapper-profile.component.scss'
})
export class TrapperProfileComponent {
  @Input() trapper!: Trapper;
}
