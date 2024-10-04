import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RanWarningComponent } from './ran-warning.component';

describe('RanWarningComponent', () => {
  let component: RanWarningComponent;
  let fixture: ComponentFixture<RanWarningComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RanWarningComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RanWarningComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
