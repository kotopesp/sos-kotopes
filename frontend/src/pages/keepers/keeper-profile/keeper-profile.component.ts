import { Component, Input } from '@angular/core';
import { FavoritesButtonComponent } from '../../../shared/post/favorites-button/favorites-button.component';
import { KeeperCostStatusComponent } from './keeper-cost-status/keeper-cost-status.component';
import { NgIf } from '@angular/common';
import { animalTypePipe } from '../../../pipes/animal_type';
import { Keeper } from '../../../model/keeper';


@Component({
  selector: 'app-keeper-profile',
  standalone: true, 
  imports: [FavoritesButtonComponent, KeeperCostStatusComponent, animalTypePipe, NgIf],
  templateUrl: './keeper-profile.component.html',
  styleUrl: './keeper-profile.component.scss'
})
export class KeeperProfileComponent {
  @Input() keeper!: Keeper;
  avatarUrl: string = '';
  ngOnInit() {
    this.avatarUrl = this.getAvatarWithBackground();
  }
  getAvatarWithBackground(): string {
    const backgrounds = ['0099ff', 'ff9900', '00cc66', 'ff66cc', '9966ff'];
    const background = backgrounds[Math.floor(Math.random() * backgrounds.length)];
    const seed = Math.random().toString(36).substring(2, 8);
    return `https://api.dicebear.com/7.x/adventurer/svg?seed=${seed}&backgroundColor=${background}`;
  }
}
