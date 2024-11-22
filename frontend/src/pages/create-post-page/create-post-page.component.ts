import {Component, ElementRef, QueryList, signal, ViewChild, ViewChildren, WritableSignal} from '@angular/core';
import {ButtonLostPetComponent} from "../../shared/buttons/button-lost-pet/button-lost-pet.component";
import {ButtonFindPetComponent} from "../../shared/buttons/button-find-pet/button-find-pet.component";
import {
  ButtonLookingForHomeComponent
} from "../../shared/buttons/button-looking-for-home/button-looking-for-home.component";
import { RouterLink} from "@angular/router";
import {DatePipe, NgClass, NgForOf, NgIf} from "@angular/common";
import {CustomCalendarComponent} from "./ui/custom-calendar/custom-calendar.component";
import {RanWarningComponent} from "../../shared/ran-warning/ran-warning.component";
import {ClickToSelectDirective} from "../../directives/click-to-select/click-to-select.directive";
import {AddPhotoButtonComponent} from "../../shared/buttons/add-photo-button/add-photo-button.component";
import {ConfirmOverlayComponent} from "../../shared/confirm-overlay/confirm-overlay.component";
import {AuthService} from "../../services/auth-service/auth.service";
import {FormsModule, ReactiveFormsModule} from "@angular/forms";
import {ChooseOneDirective} from "../../directives/choose-one/choose-one.directive";
import {PostsService} from "../../services/posts-services/posts.service";
import {ButtonStateService} from "../../services/button-state/button-state.service";

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
    NgForOf,
    CustomCalendarComponent,
    RanWarningComponent,
    DatePipe,
    ClickToSelectDirective,
    AddPhotoButtonComponent,
    ConfirmOverlayComponent,
    FormsModule,
    ChooseOneDirective,
    ReactiveFormsModule,
    NgClass,
  ],
  templateUrl: './create-post-page.component.html',
  styleUrl: './create-post-page.component.scss'
})
export class CreatePostPageComponent {
  @ViewChild('fileInput') fileInput: ElementRef | undefined;
  // ViewChild для доступа к элементу div
  @ViewChildren('myDiv') myDivs: QueryList<ElementRef> | undefined;
  titleObject: TitleObject;
  isDragging = false;
  districts: Array<{ text: string }>;
  buttonActive: boolean;
  countOfSlides: number;
  isDisabled: boolean = true;
  reason: string;
  species: string;
  gender: string;
  selectedColor: string;
  selectedFiles: { name: string, preview: string, file: File}[] = [];
  chooseColors: Array<string>;
  selectedDate!: Date;
  selectedDistrict: string;
  textValue: string;
  numberOfSlide: WritableSignal<number>;
  photosOverlay: WritableSignal<boolean>;
  descriptionOverlay: WritableSignal<boolean>;
  disableObject: { [key: number]: Array<any> };

  constructor(private authService: AuthService, private postsService: PostsService, private buttonState: ButtonStateService) {
    this.reason = '';
    this.gender = '';
    this.species = '';
    this.selectedColor = '';
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

    this.districts = [
      { text: 'м. Комендантский Проспект' },
      { text: 'р-н. Колпинский' },
    ];

    this.chooseColors = ['Чёрный', 'Белый', 'Чёрно-белый ("маркиз")', 'Полосатый', 'Рыжий', 'Серый', 'Трёхцветный']
    this.photosOverlay = signal<boolean>(false)
    this.descriptionOverlay = signal<boolean>(false)
    this.textValue = '';
    this.selectedDistrict = '';

    if (authService.Token) {
      this.countOfSlides = 6;
    } else {
      this.countOfSlides = 7
    }

    this.disableObject = {
      1: [() => this.reason],
      2: [() => this.species, () => this.gender],
      3: [() => this.selectedFiles],
      4: [() => this.selectedDistrict, () => this.selectedDate],
      5: [() => this.selectedColor],
      6: [() => this.species, () => this.gender],
      // 7: [() => this.species, () => this.gender],
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
          preview: e.target.result,
          file: e.target.file,
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
    console.log(this.selectedFiles);
  }

  createPost() {
    if (!this.textValue) {
      this.descriptionOverlay.set(true);
    }
    const form = new FormData()
    form.append('animal_type', this.species)
    this.selectedFiles.forEach((item) => {
      console.log(item);
      form.append('photo', item.file);
    });
    form.append('color', this.selectedColor)
    form.append('location', this.selectedDistrict)
    form.append('gender', this.gender)
    form.append('description', this.textValue)
    form.append('status', this.reason)

    this.postsService.createPost(form)
  }

  saveDivValue(index: number) {
    if (this.myDivs) {
      const divArray = this.myDivs.toArray();
      this.selectedDistrict = divArray[index].nativeElement.innerText;
    }
  }

  buttonNextDisabled(): boolean {
    const dependencies = this.disableObject[this.numberOfSlide()];
    if (!dependencies) return false;
    return !dependencies.every((dependencyFn) => {
      const value = dependencyFn();
      return Boolean(value);
    });
  }


  getPhotoClass(count: number): string {
    switch (count) {
      case 1:
        return '';
      case 2:
        return 'two-photos';
      case 3:
        return 'three-photos';
      case 4:
        return 'four-photos';
      default:
        return '';
    }
  }
}

