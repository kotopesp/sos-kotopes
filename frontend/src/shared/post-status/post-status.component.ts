import { Component } from '@angular/core';
import {NgClass} from "@angular/common";

@Component({
  selector: 'app-post-status',
  standalone: true,
  imports: [
    NgClass
  ],
  templateUrl: './post-status.component.html',
  styleUrl: './post-status.component.scss'
})
export class PostStatusComponent {
  className = "found";
  text = "Найден";
}
