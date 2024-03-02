<script setup lang="ts">

import {computed, PropType} from "vue";
import {Card, Core} from "@/views/Dashboard/core";
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

    if (!currentCard.value?.keysCapture) {
      return;
    }

    currentCard.value.keysCapture.forEach((act, index) => {
      if (act.keys?.has(val.keyCode)) {
        if (!act.action) {
          return
        }
        callAction(act, val.key, val.keyCode)
      }
    })
  }
})

const callAction = async (params, key, keyCode,) => {
  const {entityId, action, areaId, tags} = params;
  api.v1.interactServiceEntityCallAction({
    id: entityId,
    name: action,
    areaId: areaId,
    tags: tags,
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
