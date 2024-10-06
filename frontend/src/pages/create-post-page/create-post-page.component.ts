import {Component, ElementRef, signal, ViewChild, WritableSignal} from '@angular/core';
import {ButtonLostPetComponent} from "../../shared/buttons/button-lost-pet/button-lost-pet.component";
import {ButtonFindPetComponent} from "../../shared/buttons/button-find-pet/button-find-pet.component";
import {
  ButtonLookingForHomeComponent
} from "../../shared/buttons/button-looking-for-home/button-looking-for-home.component";
import { RouterLink} from "@angular/router";
import {DatePipe, NgForOf, NgIf, NgStyle} from "@angular/common";
import {CustomCalendarComponent} from "./ui/custom-calendar/custom-calendar.component";
import {RanWarningComponent} from "../../shared/ran-warning/ran-warning.component";
import {ClickToSelectDirective} from "../../directives/click-to-select/click-to-select.directive";
import {AddPhotoButtonComponent} from "../../shared/buttons/add-photo-button/add-photo-button.component";
import {ConfirmOverlayComponent} from "../../shared/confirm-overlay/confirm-overlay.component";
import {AuthService} from "../../services/auth-service/auth.service";
import {FormsModule} from "@angular/forms";
import {ChooseOneDirective} from "../../directives/choose-one/choose-one.directive";

interface TitleObject {
  [key: number]: string
}

@Component({
  selector: 'app-create-post-page',
  standalone: true,
  imports: [
    ButtonLostPetComponent,
    ButtonFindPetComponent,
    ButtonLookingForHomeComponent,
    RouterLink,
    NgIf,
    NgStyle,
    NgForOf,
    CustomCalendarComponent,
    RanWarningComponent,
    DatePipe,
    ClickToSelectDirective,
    AddPhotoButtonComponent,
    ConfirmOverlayComponent,
    FormsModule,
    ChooseOneDirective,
  ],
  templateUrl: './create-post-page.component.html',
  styleUrl: './create-post-page.component.scss'
})
export class CreatePostPageComponent {
  @ViewChild('fileInput') fileInput: ElementRef | undefined;
  titleObject: TitleObject;
  selectedFiles: { name: string, preview: string }[] = [];
  isDragging = false;
  selectedDate!: Date;
  selectedDistrict!: string;
  chooseColors: string[] = [];
  buttonActive: boolean;
  countOfSlides: number;

  textValue: string;
  numberOfSlide: WritableSignal<number>;
  photosOverlay: WritableSignal<boolean>;
  descriptionOverlay: WritableSignal<boolean>;

  constructor(authService: AuthService) {
    this.buttonActive = false;
    this.numberOfSlide = signal<number>(1);
    this.titleObject = {
      1: "Что случилось?",
      2: "Кто пропал?",
      3: "Прикрепите фото питомца",
      4: "Время и место",
      5: "Окрас",
      6: "Описание",
      7: "Дайте описание о вас"
    }
    this.chooseColors = ['Чёрный', 'Белый', 'Чёрно-белый ("маркиз")', 'Полосатый', 'Рыжий', 'Серый', 'Трёхцветный']
    this.photosOverlay = signal<boolean>(false)
    this.descriptionOverlay = signal<boolean>(false)
    this.textValue = '';

    if (authService.Token) {
      this.countOfSlides = 6;
    } else {
      this.countOfSlides = 7
    }
  }

  // Обработка события перетаскивания
  onDragOver(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = true;
  }

  onDragLeave(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    event.stopPropagation();
    this.isDragging = false;

    const files = event.dataTransfer?.files;
    if (files) {
      this.handleFiles(files);
    }
  }

  // Обработка выбранных файлов
  handleFiles(files: FileList): void {
    Array.from(files).forEach(file => {
      const reader = new FileReader();

      reader.onload = (e: any) => {
        this.selectedFiles.push({
          name: file.name,
          preview: e.target.result
        });
      };

      reader.readAsDataURL(file);  // Преобразование в Data URL для предварительного просмотра
    });
  }

  updateValue(value: Date) {
    this.selectedDate = value; // Обновляем переменную значением, переданным дочерним компонентом
  }

  buttonNextClick() {
    if (this.numberOfSlide() === 3 && !this.selectedFiles.length) {
      this.photosOverlay.set(true);
    } else {
      this.numberOfSlide.set(this.numberOfSlide() + 1);
    }
  }

  createPost() {
    if (!this.textValue) {
        this.descriptionOverlay.set(true);
      }
  }
}

