import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SetupOtpComponent } from './setup-otp.component';

describe('SetupOtpComponent', () => {
  let component: SetupOtpComponent;
  let fixture: ComponentFixture<SetupOtpComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SetupOtpComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SetupOtpComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
