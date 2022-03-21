import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AgentProvisioningComponent } from './agent-provisioning.component';

describe('ProvisioningComponent', () => {
  let component: AgentProvisioningComponent;
  let fixture: ComponentFixture<AgentProvisioningComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ AgentProvisioningComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AgentProvisioningComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
