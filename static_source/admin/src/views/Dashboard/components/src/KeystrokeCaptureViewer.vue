<script setup lang="ts">

import {computed, onBeforeUnmount, onMounted, PropType} from "vue";
import {Card, Core} from "@/views/Dashboard/core";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import api from "@/api/api";
import {propTypes} from "@/utils/propTypes";
import {eventBus} from "@/components/EventBus";

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
const keysHandler = (eventName: string, val: any) => {
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

onMounted(() => {
  eventBus.subscribe('keydown', keysHandler)
})

onBeforeUnmount(() => {
  eventBus.unsubscribe('keydown', keysHandler)
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
