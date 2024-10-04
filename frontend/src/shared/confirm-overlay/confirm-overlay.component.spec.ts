import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ConfirmOverlayComponent } from './confirm-overlay.component';

describe('ConfirmOverlayComponent', () => {
  let component: ConfirmOverlayComponent;
  let fixture: ComponentFixture<ConfirmOverlayComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ConfirmOverlayComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ConfirmOverlayComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
