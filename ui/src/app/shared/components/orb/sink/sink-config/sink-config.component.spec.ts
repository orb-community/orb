import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { SinkConfigComponent } from './sink-config.component';

describe('SinkConfigComponent', () => {
  let component: SinkConfigComponent;
  let fixture: ComponentFixture<SinkConfigComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ SinkConfigComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(SinkConfigComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
