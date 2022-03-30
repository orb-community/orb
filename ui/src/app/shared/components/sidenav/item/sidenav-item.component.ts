import { Component, Input, OnInit } from '@angular/core';
import { MenuItem } from '../../../interfaces/menu-item.interface';

@Component({
  selector: 'app-sidenav-item',
  templateUrl: './sidenav-item.component.html',
  styleUrls: ['./sidenav-item.component.scss']
})
export class SidenavItemComponent implements OnInit {
  @Input()
  item: MenuItem;

  constructor() {
    this.item = {};
  }

  ngOnInit(): void {
  }

}
