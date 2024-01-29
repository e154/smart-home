import {ApiEntity, ApiImage} from '@/api/stub'

export interface GridProp {
  image?: ApiImage
  key: string
  position?: boolean
  top?: number
  left?: number
  height?: number
  width?: number
}

export interface ItemPayloadGrid {
  items: GridProp[]
  image?: ApiImage
  cellHeight: number
  cellWidth: number
  showCellValue?: boolean
  gap: boolean
  gapSize: number
  attribute: string
  entity?: ApiEntity
  entityId?: string
  actionName?: string
  position?: boolean
  top?: number
  left?: number
  fontSize?: number
}
