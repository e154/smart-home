import {AttributeValue} from '@/components/Attributes'

export interface EventHTML5Notify {
  title: string
  options?: NotificationOptions
}

export interface EventTaskCompleted {
  id: number
}

export interface EventActionCompleted {
  id: number
}

export interface EventNewWebPushPublicKey {
  user_id?: number
  public_key: string
}

export interface EventTriggerCompleted {
  id: number
  args: Map<string, AttributeValue>
  entity_id: string
  last_time: string
}

export interface State {
  name: string
  description?: string
  icon?: string
  image_url?: string
}

export interface EventState {
  entity_id: string
  attributes: Map<string, AttributeValue>
  settings: Map<string, AttributeValue>
  last_changed?: string
  last_updated?: string
  state: State
  value?: string
}

export interface EventStateChange {
  storage_save: boolean
  plugin_name: string
  entity_id: string
  old_state: EventState
  new_state: EventState
}
