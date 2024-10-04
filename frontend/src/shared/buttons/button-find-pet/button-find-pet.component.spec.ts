import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ButtonFindPetComponent } from './button-find-pet.component';

describe('ButtonFindPetComponent', () => {
  let component: ButtonFindPetComponent;
  let fixture: ComponentFixture<ButtonFindPetComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ButtonFindPetComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ButtonFindPetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
