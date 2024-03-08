import { ApiImage } from '@/api/stub'
import {CompareProp} from "@/views/Dashboard/core";

export interface ImageProp extends CompareProp {
  image?: ApiImage
}

export interface ItemPayloadModal {
  items: ImageProp[]
}
