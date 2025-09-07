import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'animalType',
  standalone: true
})
export class animalTypePipe implements PipeTransform {

  transform(value: string): string {
    if (value == "cat") return "/assets/icons/cat.png"
    else if (value == "dog") return "/assets/icons/dog.png"
    else return "/assets/icons/cadog.png"
  }

}