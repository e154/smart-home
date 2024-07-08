import {ApplyFilter} from './filters'
import {AttributeValue, GetApiAttributeValue, GetAttributeValue, RenderAttributeValue} from "@/components/Attributes";
import {Cache} from "./cache";
import {ApiAttribute} from "@/api/stub";

export function Resolve(path: string, obj: any): any {
  return path.split('.').reduce(function (prev, curr) {
    return prev ? prev[curr] : null
  }, obj || self)
}

// function return array of tokens, example ['new_state.attributes.used_percent']
export function GetTokens(text?: string, cache?: Cache): string[] {
  if (!text || !cache) {
    return []
  }

  if (cache.get(text)) {
    return cache.get(text) as string[]
  }

  const regex = /(\[{1}[a-zA-Zа-яА-Я0-9_ {}|:@\-\.]{2,74}\]{1})/gm
  let tokens: string[] = text.match(regex) || []
  for (const index in tokens) {
    tokens[index] = tokens[index].replace('[', '').replace(']', '')
  }

  tokens = [...new Set(tokens)]
  cache.push(text, tokens)

  // console.log('tokens', tokens);

  return tokens || []
}

export const RenderText = async (tokens: string[], text: string, obj?: any): string => {
  let val: any

  let result = text
  for (const token of tokens) {
    const tokenFiltered = token.split('|')
    if (tokenFiltered.length > 1) {
      val = Resolve(tokenFiltered[0], obj)
    } else {
      val = Resolve(token, obj)
    }

    if (typeof val === 'object') {
      if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
        val = RenderAttributeValue(val as AttributeValue)
      }
    }

    if (tokenFiltered.length > 1) {
      val = await ApplyFilter(val, tokenFiltered[1])
    }

    if (val == undefined) {
      val = '[NO VALUE]'
    }
    // console.log(token, val)

    if (typeof val === 'object') {
      val = JSON.stringify(val)
    }

    result = result.replaceAll(`[${token}]`, val)
  }

  return result
}

export const RenderVar = (token: string, obj?: any): any => {

  // for inverse dependence
  token = token.replace('[', '').replace(']', '')

  let val: any

  const tokenFiltered = token.split('|')
  if (tokenFiltered.length > 1) {
    val = Resolve(tokenFiltered[0], obj)
  } else {
    val = Resolve(token, obj)
  }

  if (typeof val === 'object') {
    if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
      const originVal = val
      val = GetAttributeValue(originVal as AttributeValue)
      if (val == undefined) {
        val = GetApiAttributeValue(originVal as ApiAttribute)
      }
    }
  }

  if (tokenFiltered.length > 1) {
    val = ApplyFilter(val, tokenFiltered[1])
  }

  if (val == undefined) {
    val = '[NO VALUE]'
  }

  return val
}

export interface NestedObject {
  [key: string]: NestedObject | null | undefined;
}

export function getAllKeys(obj: NestedObject, parentKey = ''): string[] {
  let keys: string[] = [];

  for (const key in obj) {
    const currentKey = parentKey ? `${parentKey}.${key}` : key;

    if (typeof obj[key] === 'object' && obj[key] !== null) {
      keys.push(currentKey);
      const childKeys = getAllKeys(obj[key] as NestedObject, currentKey);
      if (childKeys.length > 0) {
        keys = keys.concat(childKeys);
      }
    } else {
      keys.push(currentKey);
    }
  }

  return keys;
}

export function getFilteredKeys(obj: NestedObject, parentKey = ''): string[] {
  let keys: string[] = [];

  for (const key in obj) {
    const currentKey = parentKey ? `${parentKey}.${key}` : key;

    if (typeof obj[key] === 'object' && obj[key] !== null) {
      if ('name' in obj[key] && 'value' in obj[key] && 'type' in obj[key]) {
        keys.push(currentKey);
      } else {
        const childKeys = getFilteredKeys(obj[key] as NestedObject, currentKey);
        if (childKeys.length > 0) {
          keys = keys.concat(childKeys);
        }
      }
    } else {
      keys.push(currentKey);
    }
  }

  return keys;
}
