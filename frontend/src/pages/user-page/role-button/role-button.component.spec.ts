import { ComponentFixture, TestBed } from '@angular/core/testing';

import { RoleButtonComponent } from './role-button.component';

describe('RoleButtonComponent', () => {
  let component: RoleButtonComponent;
  let fixture: ComponentFixture<RoleButtonComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [RoleButtonComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(RoleButtonComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
