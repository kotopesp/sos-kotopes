import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ButtonLookingForHomeComponent } from './button-looking-for-home.component';

describe('ButtonLookingForHomeComponent', () => {
  let component: ButtonLookingForHomeComponent;
  let fixture: ComponentFixture<ButtonLookingForHomeComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ButtonLookingForHomeComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ButtonLookingForHomeComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
