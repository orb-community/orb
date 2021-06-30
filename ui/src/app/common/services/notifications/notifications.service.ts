
import { Injectable } from '@angular/core';

import { NbToastrService, NbGlobalPosition, NbGlobalPhysicalPosition } from '@nebular/theme';
import { NbComponentStatus } from '@nebular/theme/components/component-status';

enum NbToastStatus {
  SUCCESS = 'success',
  INFO = 'info',
  WARNING = 'warning',
  PRIMARY = 'primary',
  DANGER = 'danger',
  DEFAULT = 'default',
}

@Injectable()
export class NotificationsService {

  private toastCfgSucc = {
    status: <NbComponentStatus>NbToastStatus.SUCCESS,
    position: <NbGlobalPosition>NbGlobalPhysicalPosition.BOTTOM_RIGHT,
  };
  private toastCfgErr = {
    status: <NbComponentStatus>NbToastStatus.DANGER,
    position: <NbGlobalPosition>NbGlobalPhysicalPosition.BOTTOM_RIGHT,
  };
  private toastCfgWarn = {
    status: <NbComponentStatus>NbToastStatus.WARNING,
    position: <NbGlobalPosition>NbGlobalPhysicalPosition.BOTTOM_RIGHT,
  };

  constructor(private toastrService: NbToastrService) { }

  success(title: string, message: string) {
    this.toastrService.show(
      message,
      title,
      this.toastCfgSucc);
  }

  error(title: string, message: string) {
    this.toastrService.show(
      message,
      title,
      this.toastCfgErr);
  }

  warn(title: string, message: string) {
    this.toastrService.show(
      message,
      title,
      this.toastCfgWarn);
  }
}
