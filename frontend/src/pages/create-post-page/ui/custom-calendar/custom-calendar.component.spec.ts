import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CustomCalendarComponent } from './custom-calendar.component';

describe('CustomCalendarComponent', () => {
  let component: CustomCalendarComponent;
  let fixture: ComponentFixture<CustomCalendarComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CustomCalendarComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CustomCalendarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
