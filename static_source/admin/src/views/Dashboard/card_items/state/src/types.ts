import {ApiImage} from '@/api/stub'
import {CompareProp} from '@/views/Dashboard/types'

export interface ImageProp extends CompareProp {
  image?: ApiImage
  icon?: string
  iconColor?: string
  iconSize?: number
}

export interface ItemPayloadState {
  items: ImageProp[]
  //deprecated
  default_image?: ApiImage
  defaultImage?: ApiImage
  defaultIcon?: string
  defaultIconColor?: string
  defaultIconSize?: number
}
