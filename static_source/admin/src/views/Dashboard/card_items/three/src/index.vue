<script setup lang="ts">
import {onMounted, PropType, ref} from "vue";
import {AmbientLight, Box, Camera, PhongMaterial, PointLight, Renderer, Scene, Texture,} from 'troisjs';
import Stats from './components/Stats'
import {CardItem} from "@/views/Dashboard/core/core";
import {GetFullUrl} from "@/utils/serverId";
import {Pane} from 'tweakpane';

// ---------------------------------
// common
// ---------------------------------


const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const el = ref<ElRef>(null)
const renderer = ref(null)
const camera = ref(null)
const box = ref(null)

const params = ref<Object>({
  color: '#ffffff',
  metalness: 1,
  roughness: 0.2,
  light1Color: '#FFFF80',
  light2Color: '#DE3578',
  light3Color: '#FF4040',
  light4Color: '#0d25bb'
})

const boxPosition = ref({y: 0, z: 0})
const planePosition = ref({y: 0, z: 0})
const cameraPosition = ref({y: 3, z: 2})

onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  // console.log(box.value.mesh.position)
  camera.value.camera.lookAt(box.value.mesh.position)

  const pane = new Pane({
    container: el.value,
  });
  pane.addBinding(params.value, 'color');
  pane.addBinding(params.value, 'metalness', {min: 0, max: 1});
  pane.addBinding(params.value, 'roughness', {min: 0, max: 1});

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

</script>

<template>
  <div ref="el" id="three-container">

    <Renderer ref="renderer" antialias resize
              :orbit-ctrl="{ autoRotate: true, enableDamping: true, dampingFactor: 0.05 }">

      <Camera ref="camera" :position="cameraPosition"/>
      <Stats :no-setup="false">

        <Scene>
          <AmbientLight :position="{ y: 50, z: 50 }"/>
          <PointLight :position="{ y: 50, z: 50 }"/>
          <!--        <Plane ref="plane" :position="planePosition"/>-->
          <Box ref="box" @pointerOver="boxOver" @click="boxClick" :position="boxPosition"
               :rotation="{ y: Math.PI / 4, z: Math.PI / 4 }">
            <PhongMaterial>
              <Texture :src="GetFullUrl('/static/assets/textures/uv-test-bw.png')"/>
            </PhongMaterial>
          </Box>
        </Scene>
      </Stats>

    </Renderer>
  </div>
</template>

<style lang="less">
.tp-rotv {
  position: absolute;
  top: 0;
  right: 0;
}
</style>
