import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePostPageComponent } from './create-post-page.component';

describe('CreatePostComponent', () => {
  let component: CreatePostPageComponent;
  let fixture: ComponentFixture<CreatePostPageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreatePostPageComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreatePostPageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
