<script setup lang="ts">

import {computed, PropType} from "vue";
import {Card, Core} from "@/views/Dashboard/core/core";
import {useEventBus} from "@/hooks/event/useEventBus";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import api from "@/api/api";
import {propTypes} from "@/utils/propTypes";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
  hover: propTypes.bool.def(false),
})

const currentCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {
  }
})

// ---------------------------------
// component methods
// ---------------------------------

useEventBus({
  name: 'keydown',
  callback: (val) => {

    if (!props.hover) {
      return
    }

    currentCard.value.keysCapture.forEach((act, index) => {
      if (act.keys?.has(val.keyCode)) {
        if (!act.entityId || !act.action) {
          return
        }
        callAction(act.entityId, act.action, val.key, val.keyCode)
      }
    })
  }
})

const callAction = async (id: string, name: string, key: string, keyCode: number) => {
  api.v1.interactServiceEntityCallAction({
    id: id,
    name: name,
    attributes: {
      "key": {
        "name": "key",
        "type": ApiTypes.STRING,
        "string": key,
      },
      "keyCode": {
        "name": "keyCode",
        "type": ApiTypes.INT,
        "int": keyCode,
      }
    },
  } as ApiEntityCallActionRequest)
}
</script>

<template>
  <div></div>
</template>

<style scoped lang="less">

</style>
