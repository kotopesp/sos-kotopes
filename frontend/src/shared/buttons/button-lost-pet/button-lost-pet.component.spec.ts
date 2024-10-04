import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ButtonLostPetComponent } from './button-lost-pet.component';

describe('ButtonLostPetComponent', () => {
  let component: ButtonLostPetComponent;
  let fixture: ComponentFixture<ButtonLostPetComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ButtonLostPetComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ButtonLostPetComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
