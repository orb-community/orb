import { ComponentFixture, TestBed } from '@angular/core/testing';

import { OneColumnLayoutComponent } from './one-column-layout.component';

describe('OneColumnLayoutComponent', () => {
  let component: OneColumnLayoutComponent;
  let fixture: ComponentFixture<OneColumnLayoutComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ OneColumnLayoutComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(OneColumnLayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
