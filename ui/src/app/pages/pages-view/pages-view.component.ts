import { Component, OnInit } from '@angular/core';
import { MENU_ITEMS } from '../menu-items';

@Component({
  selector: 'app-pages-view',
  templateUrl: './pages-view.component.html',
  styleUrls: ['./pages-view.component.scss'],
})
export class PagesViewComponent implements OnInit {
  menuItems = MENU_ITEMS;

  constructor() {
  }

  ngOnInit(): void {
  }

}
