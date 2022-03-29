import { ComponentFixture, TestBed } from '@angular/core/testing';

import { SinkListComponent } from './sink-list.component';

describe('SinksListComponent', () => {
  let component: SinkListComponent;
  let fixture: ComponentFixture<SinkListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ SinkListComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(SinkListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
