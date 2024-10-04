import {Component, Input, WritableSignal} from '@angular/core';

@Component({
  selector: 'app-add-photo-button',
  standalone: true,
  imports: [],
  templateUrl: './add-photo-button.component.html',
  styleUrl: './add-photo-button.component.scss'
})
export class AddPhotoButtonComponent {
  @Input() selectedFiles!: { name: string, preview: string }[];
  @Input() photosOverlay!: WritableSignal<boolean>

  onFileSelected(event: any): void {
    this.photosOverlay.set(false);
    const files = event.target.files;

    if (files) {
      for (let file of files) {
        const reader = new FileReader();

        reader.onload = (e: any) => {
          this.selectedFiles.push({
            name: file.name,
            preview: e.target.result
          });
        };

        reader.readAsDataURL(file);  // Чтение файла как Data URL для предварительного просмотра
      }
    }
  }
}
