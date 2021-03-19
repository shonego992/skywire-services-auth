import {Injectable} from '@angular/core';
import {ToastrService} from 'ngx-toastr';

const TOAST_TIMEOUT = 3000;

@Injectable()
export class SharedService {

  public constructor(private toastr: ToastrService) {
  }

  public sleep(ms: number) {
    return new Promise(resolve => setTimeout(resolve, ms));
  }

  public showSuccess (message: string) {
    console.info('success: ', message);
    this.toastr.success(message, '', {
      timeOut: TOAST_TIMEOUT
    });
  }

  public showError (title: string, message: string) {
    console.error('error: ', title, message);
    this.toastr.error(message, title, {
      timeOut: TOAST_TIMEOUT
    });
  }

  // Map the status from the backend, and return if the user has application in progress or can create a new one
  // Statuses from backend:
  // PENDING = 0
  // APPROVED = 1
  // DENIED_RESUBMIT = 2
  // DENIED_NO_RESUBMIT = 3
  // CANCELED = 4
  public checkIfApplicationIsInProgress(status: number): boolean {
    if (status === 0 || status === 2) {
      return true;
    }
    return false;
  }

  public mapStatus(value: number): string {
    switch (value) {
      case 0: {
        return 'PENDING';
      }
      case 1: {
        return 'APPROVED';
      }
      case 2: {
        return 'DENIED CAN RESUBMIT';
      }
      case 3: {
        return 'DENIED CANNOT RESUBMIT';
      }
      case 4: {
        return 'CANCELED';
      }
      default: {
        return 'UNKNOWN';
      }
    }
  }


}
