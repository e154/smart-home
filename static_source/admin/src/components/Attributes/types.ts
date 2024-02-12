// import Iconify from "@purge-icons/generated";
import {parseTime} from "@/utils";
import {ApiAttribute, ApiTypes} from "@/api/stub";

export enum Types {
  INT = 'int',
  STRING = 'string',
  FLOAT = 'float',
  BOOL = 'bool',
  ARRAY = 'array',
  TIME = 'time',
  MAP = 'map',
  IMAGE = 'image',
  ICON = 'icon',
  POINT = 'point',
  ENCRYPTED = 'encrypted',
}

export class EntityAttribute implements ApiAttribute {
  constructor(name: string) {
    this.name = name
    this.type = ApiTypes.STRING
    this.string = ''
  }

  name: string;
  type: ApiTypes;
  int?: number;
  string: string;
  icon: string;
  imageUrl: string;
  bool?: boolean;
  float?: number;
  array?: ApiAttribute[];
  encrypted?: string;
}

export interface AttributeValue {
  name: string
  type: Types
  value: any
}

export function GetAttributeValue(attr: AttributeValue): string {
  let val: string
  switch (attr.type.toLowerCase()) {
    case Types.INT:
      val = attr.value
      break
    case Types.STRING:
      val = attr.value
      break
    case Types.FLOAT:
      val = attr.value
      break
    case Types.BOOL:
      if (attr.value) {
        val = 'ON'
      } else {
        val = 'OFF'
      }
      break
    case Types.ARRAY:
      val = JSON.stringify(attr.value)
      break
    case Types.IMAGE:
      val = attr.value
      break
    case Types.ICON:
      // const svg = Iconify.renderHTML(attr.value, {})
      // if (svg) {
      //   val = svg
      // } else {
      val = attr.value
      // }
      break
    case Types.TIME:
      val = parseTime(attr.value) as string
      break
    case Types.MAP:
      val = JSON.stringify(attr.value)
      break
    case Types.POINT:
      val = attr.value
      break
    case Types.ENCRYPTED:
      val = attr.value
      break
    default:
      return `unknown type "${attr.type}"`
  }
  return val
}

export function GetApiAttributeValue(attr: ApiAttribute): any {
  switch (attr.type.toLowerCase()) {
    case Types.INT:
      return attr.int
    case Types.STRING:
      return attr.string
    case Types.FLOAT:
      return attr.float
    case Types.BOOL:
      if (attr.bool) {
        return 'ON'
      } else {
        return 'OFF'
      }
    case Types.ARRAY:
      return attr.array
    case Types.IMAGE:
      return attr.imageUrl
    case Types.ICON:
      // const svg = Iconify.renderHTML(attr.value, {})
      // if (svg) {
      //   return svg
      // } else {
      return attr.icon
    // }
    case Types.TIME:
      return parseTime(attr.time) as string
    case Types.MAP:
      return '[MAP]'
    case Types.POINT:
      return attr.point
    case Types.ENCRYPTED:
      return attr.encrypted
    default:
      return `unknown type "${attr.type}"`
  }
}
