import { Component } from '@angular/core';
import { FavoritesButtonComponent } from '../../../shared/post/favorites-button/favorites-button.component';
import { KeeperCostStatusComponent } from './keeper-cost-status/keeper-cost-status.component';


@Component({
  selector: 'app-keeper-profile',
  standalone: true, 
  imports: [FavoritesButtonComponent, KeeperCostStatusComponent],
  templateUrl: './keeper-profile.component.html',
  styleUrl: './keeper-profile.component.scss'
})
export class KeeperProfileComponent {

}
