import { Component } from '@angular/core';
import {ListRoleComponent} from "./ui/list-role/list-role.component";
import {NgForOf} from "@angular/common";
import {ListRole} from "../../../../model/list-role.interface";

@Component({
  selector: 'app-message-popup',
  standalone: true,
  imports: [
    ListRoleComponent,
    NgForOf
  ],
  templateUrl: './message-popup.component.html',
  styleUrl: './message-popup.component.scss'
})
export class MessagePopupComponent {
  listRoles: ListRole[] = [
    {
      color: 'brown-gray',
      icon: 'double-letter-icon.svg'
    },
    {
      color: 'light-brown',
      icon: 'white-paw-icon.svg'
    },
    {
      color: 'red-brown',
      icon: 'white-hand-icon.svg'
    },
    {
      color: 'purple',
      icon: 'white-net-icon.svg'
    },
    {
      color: 'cyan',
      icon: 'white-cross-icon.svg'
    },
    {
      color: 'dark-brown',
      icon: 'white-balls-icon.svg'
    },
  ]
}
