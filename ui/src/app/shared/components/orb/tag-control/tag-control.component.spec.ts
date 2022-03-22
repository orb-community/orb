import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { TagControlComponent } from './tag-control.component';

describe('TagControlComponent', () => {
  let component: TagControlComponent;
  let fixture: ComponentFixture<TagControlComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ TagControlComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(TagControlComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
