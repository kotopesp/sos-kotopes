import {Component, Input, WritableSignal} from '@angular/core';

@Component({
  selector: 'app-add-photo-button',
  standalone: true,
  imports: [],
  templateUrl: './add-photo-button.component.html',
  styleUrl: './add-photo-button.component.scss'
})
export class AddPhotoButtonComponent {
  @Input() selectedFiles!: { name: string, preview: string, file: File }[];
  @Input() photosOverlay!: WritableSignal<boolean>

  // Обработка выбора файлов через input
  onFileSelected(event: any): void {
    const files = event.target.files;
    if (files) {
      this.handleFiles(files);
      console.log(files)
      this.photosOverlay.set(false);
    }
  }

  // Обработка выбранных файлов
  handleFiles(files: FileList): void {

    Array.from(files).forEach(file => {
      const reader = new FileReader();

      reader.onload = (e: any) => {
        this.selectedFiles.push({
          name: file.name,
          preview: e.target.result,
          file: file
        });
      };

      reader.readAsDataURL(file);  // Преобразование в Data URL для предварительного просмотра
    });
  }
}
