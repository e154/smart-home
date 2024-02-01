import {ApiEntity, ApiImage} from '@/api/stub'

export interface TileProp {
  image?: ApiImage
  key: string
  position?: boolean
  top?: number
  left?: number
  height?: number
  width?: number
}

export interface ItemPayloadTiles {
  items: TileProp[]
  image?: ApiImage
  tileHeight: number
  tileWidth: number
  attribute: string
  entity?: ApiEntity
  entityId?: string
  actionName?: string
  position?: boolean
  top?: number
  left?: number
}
