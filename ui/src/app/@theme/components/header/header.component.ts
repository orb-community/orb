import {Component, OnDestroy, OnInit} from '@angular/core';
import {NbMediaBreakpointsService, NbMenuService, NbSidebarService, NbThemeService} from '@nebular/theme';

import {LayoutService} from '../../../@core/utils';
import {map, takeUntil} from 'rxjs/operators';
import {Observable, Subject} from 'rxjs';
import {RippleService} from '../../../@core/utils/ripple.service';

// Mainflux - Users
import {User} from 'app/common/interfaces/mainflux.interface';
import {UsersService} from 'app/common/services/users/users.service';

import {STRINGS} from 'assets/text/strings';

@Component({
  selector: 'ngx-header',
  styleUrls: ['./header.component.scss'],
  templateUrl: './header.component.html',
})
export class HeaderComponent implements OnInit, OnDestroy {
  private destroy$: Subject<void> = new Subject<void>();
  public readonly materialTheme$: Observable<boolean>;
  userPictureOnly: boolean = false;
  user: User;
  logotype: string = STRINGS.header.logotype;

  themes = [
    {
      value: 'dark',
      name: 'Dark',
    },
    {
      value: 'material-dark',
      name: 'Material Dark',
    },
  ];

  currentTheme = 'default';

  // Mainflux - Menu and version
  userMenu = [
    { title: 'Profile', link: '/pages/profile' },
    { title: 'Log out', link: '/auth/logout' },
  ];
  version = '0.0.0';

  public constructor(
    private sidebarService: NbSidebarService,
    private menuService: NbMenuService,
    private themeService: NbThemeService,
    private layoutService: LayoutService,
    private breakpointService: NbMediaBreakpointsService,
    private rippleService: RippleService,
    private usersService: UsersService,
  ) {
  }

  ngOnInit() {
    // Mfx - get Users infos and picture
    this.usersService.getProfile().subscribe(
      (resp: any) => {
        this.user = resp;
        this.user.picture = this.usersService.getUserPicture();
      },
    );
    this.usersService.getServiceVersion().subscribe(
      (resp: any) => {
        this.version = resp.version;
      },
    );

    const { xl } = this.breakpointService.getBreakpointsMap();
    this.themeService.onMediaQueryChange()
      .pipe(
        map(([, currentBreakpoint]) => currentBreakpoint.width < xl),
        takeUntil(this.destroy$),
      )
      .subscribe((isLessThanXl: boolean) => this.userPictureOnly = isLessThanXl);

    this.themeService.onThemeChange()
      .pipe(
        map(({ name }) => name),
        takeUntil(this.destroy$),
      )
      .subscribe(themeName => {
        this.currentTheme = themeName;
        this.rippleService.toggle(themeName?.startsWith('material'));
      });
  }

  ngOnDestroy() {
    this.destroy$.next();
    this.destroy$.complete();
  }

  changeTheme(themeName: string) {
    this.themeService.changeTheme(themeName);
  }

  toggleSidebar(): boolean {
    this.sidebarService.toggle(true, 'menu-sidebar');
    this.layoutService.changeLayoutSize();

    return false;
  }

  navigateHome() {
    this.menuService.navigateHome();
    return false;
  }
}
