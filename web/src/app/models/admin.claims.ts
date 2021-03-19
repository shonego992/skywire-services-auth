export class AdminClaims {
  flag_vip: boolean;
  can_create: boolean;
  can_disable: boolean;
  review_whitelist: boolean;

  constructor() {
    this.flag_vip = false;
    this.can_create = false;
    this.can_disable = false;
    this.review_whitelist = false;
  }
}
