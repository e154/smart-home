import {ApiDashboardTab, ApiEntity, ApiImage} from "@/api/stub";

export interface Area {
  id: number;
  name: string;
  description?: string;
}

export interface Dashboard {
  id: number;
  name: string;
  description?: string;
  enabled: boolean;
  areaId?: number;
  area?: Area;
  tabs?: ApiDashboardTab[];
  entities?: Map<string, ApiEntity>;
  createdAt?: string;
  updatedAt?: string;
}

// eq: равно
// lt: меньше чем
// le: меньше или равно
// ne: не равно
// ge: больше или равно
// gt: больше чем
export enum comparisonType {
  EQ = 'eq',
  LT = 'lt',
  LE = 'le',
  NE = 'ne',
  GE = 'ge',
  GT = 'gt',
}

export function Compare(x: any, y: any, rule: comparisonType): boolean {
  switch (rule) {
    case comparisonType.EQ:
      return (x == y)
    case comparisonType.LT:
      return (x < y)
    case comparisonType.LE:
      return (x <= y)
    case comparisonType.NE:
      return (x != y)
    case comparisonType.GE:
      return (x >= y)
    case comparisonType.GT:
      return (x > y)
  }
  return false
}

export interface CompareProp {
  key: string;
  comparison: comparisonType;
  value: string;
  entity?: { id?: string };
  entityId?: string;
  eventName?: string;
  eventArgs?: any;
}

export interface ButtonAction {
  entityId?: string;
  tags?: string[];
  areaId?: number;
  entity?: { id?: string };
  action: string;
  image?: ApiImage | null;
  icon?: string;
  iconColor?: string;
  iconSize?: number;
  eventName?: string;
  eventArgs?: string;
}

export interface EventContextMenu {
  event: MouseEvent;
  owner: 'card' | 'tab' | 'editor' | 'cardItem'
  tabId?: number;
  cardId?: number;
  cardItemId?: number;
}
