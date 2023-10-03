import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SinkDetailsComponent } from './sink-details.component';

describe('SinkDetailsComponent', () => {
  let component: SinkDetailsComponent;
  let fixture: ComponentFixture<SinkDetailsComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SinkDetailsComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SinkDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
