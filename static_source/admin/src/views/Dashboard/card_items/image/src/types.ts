import {ApiImage} from "@/api/stub";
import {CompareProp} from "@/views/Dashboard/core/types";

export interface ImageProp extends CompareProp {
    attrField?: string;
    image?: ApiImage
    background?: boolean
}

export interface ItemPayloadImage {
    items: ImageProp[];
    attrField?: string
    image?: ApiImage
    background?: boolean
}
