import {ApiEntity, ApiImage} from '@/api/stub'
import {CompareProp} from '@/views/dashboard/core'

export interface ImageProp extends CompareProp {
  image?: ApiImage
}

export interface ItemPayloadState {
  items: ImageProp[]
  default_image?: ApiImage
}

export interface Marker {
  image?: ApiImage
  entityId?: string
  entity?: ApiEntity
  attribute?: string
  opacity?: number
  scale?: number
}

export interface ItemPayloadMap {
  markers: Marker[]
}
