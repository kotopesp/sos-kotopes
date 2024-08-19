import {Component, Input} from '@angular/core';
import {ListRole} from "../../../../../../model/list-role.interface";

@Component({
  selector: 'app-list-role',
  standalone: true,
  imports: [],
  templateUrl: './list-role.component.html',
  styleUrl: './list-role.component.scss'
})
export class ListRoleComponent {

  @Input() listRole: ListRole;

  constructor() {
    this.listRole = new class ListRole implements ListRole {
      color: string;
      icon: string;

      constructor() {
        this.color = '';
        this.icon = '';
      }
    }
  }
}
