import { FavoritesButtonComponent } from '../../../shared/post/favorites-button/favorites-button.component';
import { FindStatusFlagComponent } from "../../../shared/post/find-status-flag/find-status-flag.component";
import { Seeker } from '../../../model/seeker';
import { SeekerTypePipe } from "../../../pipes/seeker-type.pipe";
import { SeekerCostStatusComponent } from "./seeker-cost-status/seeker-cost-status.component";
import { Component, Input } from '@angular/core';
import { NgIf } from '@angular/common';

@Component({  
  selector: 'app-seeker-profile',
  standalone: true,
  imports: [FavoritesButtonComponent, FindStatusFlagComponent, SeekerTypePipe, SeekerCostStatusComponent, NgIf],
  templateUrl: './seeker-profile.component.html',
  styleUrl: './seeker-profile.component.scss'
})
export class SeekerProfileComponent {
  @Input() seeker!: Seeker;
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
