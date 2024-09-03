import {Component, Input, signal, WritableSignal} from '@angular/core';
import {WriteButtonComponent} from "./ui/write-button/write-button.component";
import {NgForOf} from "@angular/common";

@Component({
  selector: 'app-write-overlay',
  standalone: true,
  imports: [
    WriteButtonComponent,
    NgForOf
  ],
  templateUrl: './write-overlay.component.html',
  styleUrl: './write-overlay.component.scss'
})
export class WriteOverlayComponent {

  @Input() writeOverlay: WritableSignal<boolean>;

  WriteButtons = [
    {
      buttonColor: 'var(--role-orange-color)',
      icon: 'white-hand-icon.svg',
      title: 'Написать как передержке'
    },
    {
      buttonColor: 'var(--role-purple-color)',
      icon: 'white-net-icon.svg',
      title: 'Написать как отловщику'
    },
    {
      buttonColor: 'var(--role-green-color)',
      icon: 'white-cross-icon.svg',
      title: 'Написать как ветеринару'
    },
    {
      buttonColor: 'var(--role-orange-color)',
      icon: 'white-balls-icon.svg',
      title: 'Написать с другой целью'
    }
  ]

  constructor() {
    this.writeOverlay = signal<boolean>(false);
  }
}
