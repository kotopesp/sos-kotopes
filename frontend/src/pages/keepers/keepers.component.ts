import { Component } from '@angular/core';
import { KeeperProfileComponent } from './keeper-profile/keeper-profile.component';
import { FavoritesButtonComponent } from '../../shared/post/favorites-button/favorites-button.component';
import { KeepersFilterBarComponent } from './keepers-filter-bar/keepers-filter-bar.component';

@Component({
  selector: 'app-keepers',
  standalone: true,
  imports: [KeeperProfileComponent, FavoritesButtonComponent, KeepersFilterBarComponent],
  templateUrl: './keepers.component.html',
  styleUrl: './keepers.component.scss'
})
export class KeepersComponent {

}
