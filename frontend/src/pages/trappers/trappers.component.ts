
import { TrapperProfileComponent } from "./trapper-profile/trapper-profile.component";
import { TrapperFilterBarComponent } from "./trapper-filter-bar/trapper-filter-bar.component";
import { Meta, Trapper } from '../../model/trapper';
import { TrapperService } from '../../services/trapper-services/trapper.service';
import { Component, } from "@angular/core";
import { TrapperFilterService } from "../../services/trapper-services/trapper-filter-service/trapper-filter.service";


@Component({
  selector: 'app-trappers',
  standalone: true,
  imports: [TrapperProfileComponent, TrapperFilterBarComponent],
  templateUrl: './trappers.component.html',
  styleUrl: './trappers.component.scss'
})
export class TrappersComponent{
  trappers: Trapper[] = [];
  meta: Meta | null = null;
  constructor(private trapperService: TrapperService, private filterService: TrapperFilterService) {
   this.updateTrappersArray()
  }
  updateTrappersArray(){
    this.trapperService.getTrappersProfile().subscribe(response => {
      if (response) {
        this.trappers = response.data.trappers;  
        this.meta = response.data.meta;
      }
    });
  }

  ngOnInit(): void {
    this.filterService.filterTags.subscribe(this.updateTrappersArray);
    this.filterService.location.subscribe(this.updateTrappersArray);
  }
}
