import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbTimeInputComponent } from './ktb-time-input.component';
import { AppModule } from '../../app.module';

describe('KtbTimeInputComponent', () => {
  let component: KtbTimeInputComponent;
  let fixture: ComponentFixture<KtbTimeInputComponent>;

  const formControlNames = ['hours', 'minutes', 'seconds', 'millis', 'micros'];

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbTimeInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should validate input for formControls with a min value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(-1);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeForm.controls[control].value).toEqual(0);
    }
  });

  it('should validate input for formControls with a max value and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(25);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeForm.controls[control].value).toEqual(24);
    }
  });

  it('should validate input for formControls, round input and set appropriate value to formControl', () => {
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(1.25);

      // when
      component.validateInput(control, 0, 24);

      // then
      expect(component.timeForm.controls[control].value).toEqual(1);
    }
  });

  it('should emit given values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(1);

      // when
      component.validateInput(control, 0, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({ hours: 1, minutes: 1, seconds: 1, millis: 1, micros: 1 });
  });

  it('should emit 0 as values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(0);

      // when
      component.validateInput(control, 0, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({ hours: 0, minutes: 0, seconds: 0, millis: 0, micros: 0 });
  });

  it('should emit undefined for not given values', () => {
    const spy = jest.spyOn(component.timeChanged, 'emit');
    for (const control of formControlNames) {
      // given
      component.timeForm.controls[control].setValue(null);

      // when
      component.validateInput(control, 0, 24);
    }

    // then
    expect(spy).toHaveBeenCalledWith({
      hours: undefined,
      minutes: undefined,
      seconds: undefined,
      millis: undefined,
      micros: undefined,
    });
  });
});
