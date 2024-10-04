import { ComponentFixture, TestBed } from '@angular/core/testing';

import { AddPhotoButtonComponent } from './add-photo-button.component';

describe('AddPhotoButtonComponent', () => {
  let component: AddPhotoButtonComponent;
  let fixture: ComponentFixture<AddPhotoButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AddPhotoButtonComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(AddPhotoButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
