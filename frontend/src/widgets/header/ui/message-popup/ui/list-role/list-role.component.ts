import {Component, Input} from '@angular/core';

@Component({
  selector: 'app-list-role',
  standalone: true,
  imports: [],
  templateUrl: './list-role.component.html',
  styleUrl: './list-role.component.scss'
})
export class ListRoleComponent {

  @Input() listRole: any;

}
