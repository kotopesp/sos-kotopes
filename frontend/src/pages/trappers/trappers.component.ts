
import { TrapperProfileComponent } from "./trapper-profile/trapper-profile.component";
import { TrapperFilterBarComponent } from "./trapper-filter-bar/trapper-filter-bar.component";
import { Trapper } from '../../model/trapper';
import { TrapperService } from '../../services/trapper-services/trapper.service';
import { Component, inject } from "@angular/core";


@Component({
  selector: 'app-trappers',
  standalone: true,
  imports: [TrapperProfileComponent, TrapperFilterBarComponent],
  templateUrl: './trappers.component.html',
  styleUrl: './trappers.component.scss'
})
export class TrappersComponent{
  trappersService = inject(TrapperService)
  trappers = this.trappersService.trappers
}
