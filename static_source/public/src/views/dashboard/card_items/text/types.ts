import { CompareProp } from '@/views/dashboard/core'

export interface TextProp extends CompareProp {
  text?: string
  tokens?: string[]
}

export interface ItemPayloadText {
  items: TextProp[]
  default_text?: string
  current_text: string
}
