import { ApiImage } from '@/api/stub'
import {CompareProp} from "@/views/Dashboard/types";

export interface ImageProp extends CompareProp {
  image?: ApiImage
}

export interface ItemPayloadDummy {
  items: ImageProp[]
  //deprecated
  default_image?: ApiImage
  defaultImage?: ApiImage
}
