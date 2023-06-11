import { Attribute, EventStateChange, GetAttrValue } from '@/api/stream_types'
import { comparisonType } from '@/views/dashboard/core'
import { ApplyFilter } from '@/views/dashboard/filters'

export function Compare(x: any, y: any, rule: comparisonType): boolean {
  switch (rule) {
    case comparisonType.EQ:
      return (x == y)
    case comparisonType.LT:
      return (x < y)
    case comparisonType.LE:
      return (x <= y)
    case comparisonType.NE:
      return (x != y)
    case comparisonType.GE:
      return (x >= y)
    case comparisonType.GT:
      return (x > y)
  }
  return false
}

export function Resolve(path: string, obj: any): any {
  return path.split('.').reduce(function(prev, curr) {
    return prev ? prev[curr] : null
  }, obj || self)
}

export class Cache {
  private pull: object;

  constructor() {
    this.pull = {}
  }

  push(key: string, value: any) {
    this.pull[key] = value
  }

  get(key: string): any | null {
    if (!this.pull.hasOwnProperty(key)) {
      return null
    }
    return this.pull[key]
  }

  clear() {
    this.pull = {}
  }
}

export function GetTokens(text?: string, cache?: Cache): string[] {
  if (!text || !cache) {
    return []
  }

  if (cache.get(text)) {
    return cache.get(text) as string[]
  }

  const regex = /(\[{1}[a-zA-Zа-яА-Я0-9_ {}|:@\-\.]{2,64}\]{1})/gm
  let tokens: string[] = text.match(regex) || []
  for (const index in tokens) {
    tokens[index] = tokens[index].replace('[', '').replace(']', '')
  }

  tokens = [...new Set(tokens)]
  cache.push(text, tokens)

  // console.log('tokens', tokens);

  return tokens || []
}

export function RenderText(tokens: string[], text: string, lastEvent?: EventStateChange): string {
  let val: any

  let result = text
  for (const token of tokens) {
    const tokenFiltered = token.split('|')
    if (tokenFiltered.length > 1) {
      val = Resolve(tokenFiltered[0], lastEvent)
    } else {
      val = Resolve(token, lastEvent)
    }

    if (typeof val === 'object') {
      if (val && val.hasOwnProperty('type') && val.hasOwnProperty('name')) {
        val = GetAttrValue(val as Attribute)
      }
    }

    if (tokenFiltered.length > 1) {
      val = ApplyFilter(val, tokenFiltered[1])
    }

    if (val == undefined) {
      val = '[NO VALUE]'
    }
    // console.log(token, key, val)

    result = result.replaceAll(`[${token}]`, val)
  }

  return result
}
