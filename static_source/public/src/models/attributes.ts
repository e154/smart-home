import { ApiAttribute, ApiTypes } from '@/api/stub'

export class Attribute implements ApiAttribute {
  constructor(name: string) {
    this.name = name
    this.type = ApiTypes.STRING
    this.string = ''
  }

  name: string;
  type: ApiTypes;
  int?: number;
  string: string;
  bool?: boolean;
  float?: number;
  array?: ApiAttribute[];
}
