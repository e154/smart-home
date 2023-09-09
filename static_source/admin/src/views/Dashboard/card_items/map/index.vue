<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch, inject} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import {ElImage, ElIcon} from "element-plus";
import { GeoJSON } from "ol/format"
import { Fill, Stroke, Style } from "ol/style"
import type { ObjectEvent } from "ol/Object";
import type { View } from "ol";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {debounce} from "lodash-es";
import markerIcon from "@/assets/imgs/marker.png";
import {ApiImage} from "@/api/stub";

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
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)
})

const center = ref([34, 39.13]);
const projection = ref("EPSG:4326");
const zoom = ref(17);
const rotation = ref(0);

const overrideStyleFunction = (feature, style, resolution) => {
  console.log({ feature, style, resolution });
  const clusteredFeatures = feature.get("features");
  const size = clusteredFeatures.length;
  style.getText().setText(size.toString());
};

const getRandomInRange = (from, to, fixed) => {
  return (Math.random() * (to - from) + from).toFixed(fixed) * 1;
};

const view = ref<View>();
const position = ref([]);
const geoLocChange = (event: ObjectEvent) => {
  // console.log("AAAAA", event);
  position.value = event.target.getPosition();
  view.value?.setCenter(event.target?.getPosition());
};

const overviewmapcontrol = ref(true);
const fullscreencontrol = ref(true);
const mousepositioncontrol = ref(true);
const attributioncontrol = ref(true);
const vectorsource = ref(null);

// ---------------------------------
// component methods
// ---------------------------------

const Feature = inject("ol-feature");
const Geom = inject("ol-geom");

// ---------------------------------
// markers methods
// ---------------------------------

const _cache = new Cache()
const update = () => {
  for (let index in props.item?.payload.map?.markers) {
    let v: string = props.item?.payload.map?.markers[index].attribute || ''
    const tokens = GetTokens(props.item?.payload.map?.markers[index].attribute, _cache)
    if (tokens.length) {
      v = RenderText(tokens, v, props.item?.lastEvent)
    }
    console.log(v)
  }
}


onMounted(() => {
  setTimeout(() => {
    for (let index in props.item?.payload.map?.markers) {
      requestCurrentState(props.item?.payload.map?.markers[index].entityId);
    }
  }, 1000);
})

const getUrl = (image?: ApiImage): string | undefined => {
  return import.meta.env.VITE_API_BASEPATH as string + image?.url || undefined;
}

update()


</script>

<template>
  <div ref="el" class="w-[100%] h-[100%] overflow-hidden">

    <ol-map
        :loadTilesWhileAnimating="true"
        :loadTilesWhileInteracting="true"
        style="height: 100%"
        ref="map"
    >
      <ol-view
          ref="view"
          :center="center"
          :rotation="rotation"
          :zoom="zoom"
          :projection="projection"
      />

      <ol-fullscreen-control v-if="fullscreencontrol" />
      <ol-mouseposition-control v-if="mousepositioncontrol"/>
      <ol-attribution-control v-if="attributioncontrol" />

      <ol-tile-layer>
        <ol-source-osm />
      </ol-tile-layer>


      <ol-vector-layer
          :updateWhileAnimating="true"
          :updateWhileInteracting="true"
      >
        <ol-source-vector ref="vectorsource">
          <ol-animation-fade :duration="4000">
<!--            <ol-feature v-for="index in 20" :key="index">-->
            <ol-feature v-for="(marker, index) in props.item?.payload.map?.markers" :key="index">
              <ol-geom-point
                  :coordinates="[
                getRandomInRange(24, 45, 3),
                getRandomInRange(35, 41, 3),
              ]"
              />

              <ol-style>
                <ol-style-icon :src="getUrl(marker.image) || markerIcon" :opacity="marker?.opacity || 0.9" :scale="marker?.scale || 0.08" />
              </ol-style>
            </ol-feature>
          </ol-animation-fade>
        </ol-source-vector>
      </ol-vector-layer>

      <ol-geolocation :projection="projection" @change:position="geoLocChange">
        <template>
          <ol-vector-layer :zIndex="2">
            <ol-source-vector>
              <ol-feature ref="positionFeature">
                <ol-geom-point :coordinates="position"/>
                <ol-style>
                  <ol-style-icon :src="markerIcon" :scale="0.02"/>
                </ol-style>
              </ol-feature>
            </ol-source-vector>
          </ol-vector-layer>
        </template>
      </ol-geolocation>
    </ol-map>

  </div>
</template>

<style lang="less" >

</style>
