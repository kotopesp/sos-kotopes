import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
  name: 'trapperType',
  standalone: true
})
export class TrapperTypePipe implements PipeTransform {

  transform(value: string): string {
    if (value == "cats") return "../assets/icons/cat.png"
    else if (value == "dogs") return "../assets/icons/dog.png"
    else return "../assets/icons/cadog.png"
  }

}
