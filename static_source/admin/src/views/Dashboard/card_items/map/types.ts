import {ApiEntity, ApiImage} from '@/api/stub'
import {CompareProp} from '@/views/dashboard/core'

export interface ImageProp extends CompareProp {
  image?: ApiImage
}

export interface ItemPayloadState {
  items: ImageProp[]
  //depredcated
  default_image?: ApiImage
  defaultImage?: ApiImage
}

export interface Marker {
  image?: ApiImage
  entityId?: string
  entity?: ApiEntity
  attribute?: string
  opacity?: number
  scale?: number
  value?: number[]
}

export interface ItemPayloadMap {
  projection: string
  zoom: number
  rotation: number
  resolution: number
  center: number[]
  staticCenter: boolean
  indexMarkerCenter?: number
  markers: Marker[]
}
