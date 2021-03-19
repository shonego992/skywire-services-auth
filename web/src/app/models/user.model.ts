import {WhitelistApplication} from './application-model';

export interface User {
  id: number;
  status: number;
  username: string;
  disabled: Date;
  password: string;
  skycoinAddress?: string;
  applications: WhitelistApplication[];
  rights: Right[];
  useOtp?: boolean;
}

export interface Right {
  Name: string;
  Label: string;
  Value: boolean;
}
