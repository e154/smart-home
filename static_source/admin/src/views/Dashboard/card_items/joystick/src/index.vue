<script setup lang="ts">
import {onMounted, PropType, ref, watch} from "vue";
import {ElIcon, ElImage} from "element-plus";
import {CardItem} from "@/views/Dashboard/core";
import {JoystickAction, JoystickController, point} from "./types";
import {useEmitt} from "@/hooks/web/useEmitt";
import {Picture as IconPicture} from '@element-plus/icons-vue'
import {debounce} from "lodash-es";
import api from "@/api/api";
import {ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import {UUID} from "uuid-generator-ts";
import {GetFullUrl} from "@/utils/serverId";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const joystick = ref()
const stick = ref()
const currentID = ref('')

const el = ref<ElRef>(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  const uuid = new UUID()
  currentID.value = uuid.getDashFreeUUID()

  joystick.value = new JoystickController(stick, 64, 8, currentID.value)
})

// ---------------------------------
// component methods
// ---------------------------------

const callAction = async (action: JoystickAction, val: point) => {
  await api.v1.interactServiceEntityCallAction({
    id: props.item?.entityId,
    name: action.action,
    tags: action.tags,
    areaId: action.areaID,
    attributes: {
      "X": {
        "name": "X",
        "type": ApiTypes.FLOAT,
        "float": val.x,
      },
      "Y": {
        "name": "Y",
        "type": ApiTypes.FLOAT,
        "float": val.y,
      }
    },
  } as ApiEntityCallActionRequest);
}

const isStarted = ref(false)
const action = debounce((val: point) => {
  // console.log(val)
  if (val.x === 0 && val.y === 0) {
    // end action
    if (isStarted.value) {
      isStarted.value = false
      if (props.item?.payload?.joystick?.endAction) {
        callAction(props.item.payload.joystick.endAction, val)
      }
    }
  } else {
    // start action
    if (!isStarted.value) {
      isStarted.value = true
      if (props.item?.payload?.joystick?.startAction) {
        callAction(props.item.payload.joystick.startAction, val)
      }
    }
  }
}, 100)

useEmitt({
  name: 'updateValue',
  callback: (val) => {
    const {id, value} = val
    if (id !== currentID.value) {
      return
    }
    action(value)
  }
})

watch(
    () => props.item?.entityId,
    (val?: string) => {

    },
)

const getStickUrl = (): string => {
  return GetFullUrl(props.item?.payload.joystick?.stickImage?.url);
}

const loop = () => {
  requestAnimationFrame(loop);
}

loop();

</script>

<template>
  <div ref="el" :class="[{'hidden': item.hidden}]">
    <div style="position: absolute; left:0; top:0; cursor: pointer;" ref="stick">
      <ElImage :src="getStickUrl()">
        <template #error>
          <div class="image-slot">
            <ElIcon>
              <icon-picture/>
            </ElIcon>
          </div>
        </template>
      </ElImage>
    </div>
  </div>
</template>

<style lang="less">

</style>
