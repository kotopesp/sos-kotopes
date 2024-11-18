import { TestBed } from '@angular/core/testing';

import { ButtonStateService } from './button-state.service';

describe('ButtonStateService', () => {
  let service: ButtonStateService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(ButtonStateService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
