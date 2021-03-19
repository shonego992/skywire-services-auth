export class WhitelistApplication {
  id: number;
  createdAt: string;
  currentStatus: any;
  userId?: number;
  changeHistory?: any[];

  constructor() {
    this.id = null;
    this.currentStatus = null;
    this.changeHistory = [];
    this.createdAt = '';
    this.userId = null;
  }
}
