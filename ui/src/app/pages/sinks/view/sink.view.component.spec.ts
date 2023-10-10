import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SinkViewComponent } from './sink.view.component';

describe('SinkViewComponent', () => {
  let component: SinkViewComponent;
  let fixture: ComponentFixture<SinkViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SinkViewComponent ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SinkViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
