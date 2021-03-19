export class ChangeApplicationStatusReq {
  applicationId: number;
  status: number;
  comment: string;

  constructor() {
    this.applicationId = null;
    this.status = null;
    this.comment = '';
  }
}
