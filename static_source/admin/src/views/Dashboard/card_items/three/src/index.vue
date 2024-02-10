<script setup lang="ts">
import {onMounted, PropType, ref} from "vue";
import {
  Box,
  Plane,
  Camera,
  PhongMaterial,
  Texture,
  Renderer,
  Scene,
  AmbientLight,
  PointLight,
} from 'troisjs';
import {Vector3} from 'three';
import {CardItem} from "@/views/Dashboard/core/core";
import {RenderVar} from "@/views/Dashboard/core/render";
import {prepareUrl} from "@/utils/serverId";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref(null)
const camera = ref(null)
const box = ref(null)

onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  // console.log(box.value.mesh.position)
  camera.value.camera.lookAt(box.value.mesh.position)

  // console.log(camera.value.camera)
})

// ---------------------------------
// component methods
// ---------------------------------

const boxColor = ref('#ffffff')

const boxOver = ({over}) => {
  boxColor.value = over ? '#ff0000' : '#ffffff';
}

const boxClick = (e) => {
  console.log(e);
}

const boxPosition = ref({y: 3, z: 0})
const planePosition = ref({y: 0, z: 0})
const cameraPosition = ref({y: 3, z: 2})

const getUrl = (url: string): string => {
  return prepareUrl(import.meta.env.VITE_API_BASEPATH as string + url);
}


</script>

<template>
  <div ref="el">

    <Renderer ref="renderer" antialias resize
              :orbit-ctrl="{ autoRotate: true, enableDamping: true, dampingFactor: 0.05 }">
      <Camera ref="camera" :position="cameraPosition"/>
      <Scene>
        <AmbientLight :position="{ y: 50, z: 50 }"/>
        <PointLight :position="{ y: 50, z: 50 }"/>
        <Plane ref="plane" :position="planePosition"/>
        <Box ref="box" @pointerOver="boxOver" @click="boxClick" :position="boxPosition"
             :rotation="{ y: Math.PI / 4, z: Math.PI / 4 }">
          <PhongMaterial>
            <Texture :src="getUrl('/static/assets/textures/uv-test-bw.png')"/>
          </PhongMaterial>
        </Box>
      </Scene>
    </Renderer>
  </div>
</template>

<style lang="less">

</style>
