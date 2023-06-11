import { ApiNewEntityRequest, ApiEntity } from '@/api/stub'
import api from '@/api/api'
import router from '@/router'

export const createCard = async(val: ApiNewEntityRequest) => {
  const entity: ApiNewEntityRequest = {
    name: val.name,
    pluginName: val.pluginName,
    description: val.description,
    area: val.area,
    icon: val.icon,
    image: val.image,
    autoLoad: val.autoLoad,
    parent: val.parent || undefined,
    actions: [],
    states: [],
    attributes: val.attributes,
    settings: val.settings,
    scripts: val.scripts
  }

  // update image
  if (entity.image) {
    entity.image = { id: entity.image.id }
  }

  // update actions
  for (const i in val.actions) {
    const action = Object.assign({}, val.actions[<any> i])
    if (action.image?.id) {
      action.image = { id: action.image?.id }
    }
    if (action.script?.id) {
      action.script = { id: action.script?.id }
    }
    entity.actions?.push(action)
  }

  // update states
  for (const i in val.states) {
    const state = Object.assign({}, val.states[<any> i])
    if (state.image?.id) {
      state.image = { id: state.image?.id }
    }
    entity.states?.push(state)
  }

  return api.v1.entityServiceAddEntity(entity)
}

export const importEntity = async(val: ApiEntity) => {
  const entity: ApiEntity = {
    id: val.id,
    pluginName: val.pluginName,
    description: val.description,
    area: val.area,
    icon: val.icon,
    image: val.image,
    autoLoad: val.autoLoad,
    parent: val.parent || undefined,
    actions: [],
    states: [],
    attributes: val.attributes,
    settings: val.settings,
    scripts: val.scripts
  }

  // update actions
  for (const i in val.actions) {
    const action = Object.assign({}, val.actions[<any> i])
    entity.actions?.push(action)
  }

  // update states
  for (const i in val.states) {
    const state = Object.assign({}, val.states[<any> i])
    entity.states?.push(state)
  }

  return api.v1.entityServiceImportEntity(entity)
}
