import {WhitelistNodeModel} from './whitelist-nodes-model';
import {WhitelistImageModel} from './whitelist-image-model';

export class ChangeHistory {
  description?: string;
  id: number;
  images: WhitelistImageModel [];
  location?: string;
  nodes: WhitelistNodeModel [];
  status: number;
  createdAt?: string;

  constructor() {
    this.description = '';
    this.id = null;
    this.images = [];
    this.location = '';
    this.nodes = [];
    this.status = null;
    this.createdAt = '';
  }
}
