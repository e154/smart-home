import { ApiImage } from '@/api/stub'
import { CompareProp } from '@/views/dashboard/core'

export interface ImageProp extends CompareProp {
  image?: ApiImage
}

export interface ItemPayloadState {
  items: ImageProp[]
  default_image?: ApiImage
}
