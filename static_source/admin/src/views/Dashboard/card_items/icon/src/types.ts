import {CompareProp} from "@/views/Dashboard/core";

export interface IconProp extends CompareProp {
  icon?: string
  attrField?: string;
  iconColor?: string
}

export interface ItemPayloadIcon {
  items: IconProp[];
  attrField?: string;
  value?: string;
  iconColor?: string;
}
