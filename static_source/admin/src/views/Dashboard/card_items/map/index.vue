<script setup lang="ts">
import {computed, inject, onMounted, PropType, ref, unref, watch} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import type {ObjectEvent} from "ol/Object";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import markerIcon from "@/assets/imgs/marker.png";
import {ApiImage} from "@/api/stub";
import {debounce} from "lodash-es";
import {View} from "ol";
import {propTypes} from "@/utils/propTypes";
import {useAppStore} from "@/store/modules/app";

// ---------------------------------
// common
// ---------------------------------

const appStore = useAppStore()

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
  editor: propTypes.bool.def(false),
})

const currentCardItem = computed(() => props.item as CardItem)

// ---------------------------------
// component methods
// ---------------------------------

const loaded = ref(false)

const center = ref(props.item?.payload.map?.center || [34, 39.13]);
const projection = ref("EPSG:4326");
const zoom = ref(props.item?.payload.map?.zoom || 6);
const rotation = ref(0);
const jawgLayer = ref(null);
const osmLayer = ref(null);

const manualMode = ref(false);
const fullscreencontrol = ref(true);
const overviewmapcontrol = ref(true);
const mousepositioncontrol = ref(false);
const attributioncontrol = ref(false);
const Feature = inject("ol-feature");
const Geom = inject("ol-geom");

const view = ref<View>();
const position = ref([]);
const geoLocChange = (event: ObjectEvent) => {
  // console.log("AAAAA", event);
  // position.value = event.target.getPosition();
  // view.value?.setCenter(event.target?.getPosition());
};

const reloadKey = ref(0)
const reload = debounce(() => {
      reloadKey.value += 1
    }, 100
)

const el = ref(null)
onMounted(() => {
  // store dom element moveable
  props.item.setTarget(el.value)

  layerList.value.push(jawgLayer.value.tileLayer);
  layerList.value.push(osmLayer.value.tileLayer);
})

// ---------------------------------
// view
// ---------------------------------
const layerList = ref([]);
const showSwipeControl = ref(true);

const compareArrays = (a, b) => {
  return JSON.stringify(a) === JSON.stringify(b);
};

const updateZoom = debounce((z) => {
  currentCardItem.value.payload.map.zoom = z
}, 100)
function zoomChanged(z) {
  if (currentCardItem.value.payload.map.zoom === z ) {
    return
  }
  // manualMode.value = !props.editor
  updateZoom(unref(z))
}

function resolutionChanged(r) {
  // currentCardItem.value.payload.map.resolution = r
}

const updateCenter = debounce((c) => {
  currentCardItem.value.payload.map.center = c
}, 100)

function centerChanged(c) {
  if (compareArrays(unref(currentCardItem.value.payload.map.center), unref(c))) {
    return
  }
  // manualMode.value = !props.editor
  updateCenter(unref(c))
}

function rotationChanged(r) {
  // currentCardItem.value.payload.map.rotation = r
}

// ---------------------------------
// markers methods
// ---------------------------------

export interface Marker {
  image?: ApiImage
  opacity?: number
  scale?: number
  position?: number[]
}

const markers = ref<Marker[]>([])

const _cache = new Cache()
const update = debounce(() => {
  console.log('update')

  loaded.value = false
  markers.value = []

  if (
      !manualMode.value && currentCardItem.value.payload.map?.staticCenter &&
      !compareArrays(unref(center.value), unref(currentCardItem.value.payload.map.center))
  ) {
    center.value = currentCardItem.value?.payload.map?.center || [0, 0]
  }

  if ( !manualMode.value && unref(zoom.value) != unref(currentCardItem.value?.payload.map.zoom)) {
    zoom.value = currentCardItem.value?.payload.map.zoom
  }

  for (let index in currentCardItem.value?.payload.map?.markers) {
    const entityId = currentCardItem.value?.payload.map?.markers[index]?.entityId;
    if (!entityId || !currentCardItem.value.payload.map?.markers[index]) {
      loaded.value = true
      return
    }
    let v: string = currentCardItem.value?.payload.map?.markers[index].attribute || ''
    const tokens = GetTokens(v, _cache)
    if (tokens.length) {
      const lastState = currentCardItem.value?.lastEvents(entityId);
      v = RenderText(tokens, v, lastState)
      if (v !== '[NO VALUE]') {
        const strArr = v.split(',');
        if (strArr.length === 2) {
          const marker = {
            image: currentCardItem.value?.payload.map?.markers[index].image,
            opacity: currentCardItem.value?.payload.map?.markers[index].opacity,
            scale: currentCardItem.value?.payload.map?.markers[index].scale,
            position: [
              parseFloat(strArr[0]),
              parseFloat(strArr[1]),
            ]
          } as Marker
          markers.value.push(marker)
          if (
              !manualMode.value && !currentCardItem.value.payload.map?.staticCenter &&
              currentCardItem.value.payload.map?.indexMarkerCenter == index
          ) {
            if (!compareArrays(unref(center.value), unref(marker.position))) {
              center.value = marker.position
            }

          }
        }
      }
    }
  }
  loaded.value = true
    }, 100
)

watch(
    () => props.item,
    (val?) => {
      if (!val) return;
      update()
    },
    {
      deep: true
    }
)

const getUrl = (image?: ApiImage): string | undefined => {
  return import.meta.env.VITE_API_BASEPATH as string + image?.url || undefined;
}

</script>

<template>
  <div ref="el" class="w-[100%] h-[100%] overflow-hidden">
    <ol-map
        :loadTilesWhileAnimating="true"
        :loadTilesWhileInteracting="true"
        :pixelRatio="1"
        style="height: 100%"
        ref="map"
        v-show="loaded"
    >
      <ol-view
          ref="view"
          :center="center"
          :rotation="rotation"
          :zoom="zoom"
          :projection="projection"
          @zoomChanged="zoomChanged"
          @centerChanged="centerChanged"
          @resolutionChanged="resolutionChanged"
          @rotationChanged="rotationChanged"
      />

      <ol-mouseposition-control v-if="mousepositioncontrol"/>
      <ol-fullscreen-control v-if="fullscreencontrol"/>
      <ol-attribution-control v-if="attributioncontrol"/>

      <ol-overviewmap-control v-if="overviewmapcontrol">
        <ol-tile-layer>
          <ol-source-osm />
        </ol-tile-layer>
      </ol-overviewmap-control>


      <ol-swipe-control
          ref="swipeControl"
          v-if="appStore.isDark && layerList.length > 0"
          :position="1.1"
          :layerList="layerList"
      />

      <ol-tile-layer ref="jawgLayer" title="JAWG">
        <ol-source-xyz
            crossOrigin="anonymous"
            url="https://c.tile.jawg.io/jawg-dark/{z}/{x}/{y}.png?access-token=87PWIbRaZAGNmYDjlYsLkeTVJpQeCfl2Y61mcHopxXqSdxXExoTLEv7dwqBwSWuJ"
        />
      </ol-tile-layer>

      <ol-tile-layer ref="osmLayer">
        <ol-source-osm />
      </ol-tile-layer>

      <ol-vector-layer
          :updateWhileAnimating="true"
          :updateWhileInteracting="true"
      >
        <ol-source-vector ref="vectorsource">
          <ol-animation-fade :duration="4000" :key="reloadKey">
            <ol-feature v-for="(marker, index) in markers" :key="index">
              <ol-geom-point :coordinates="marker.position"/>
              <ol-style>
                <ol-style-icon
                    :src="getUrl(marker.image || markerIcon)"
                    :opacity="marker?.opacity || 0.9"
                    :scale="marker?.scale || 0.08"/>
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

<style lang="less">
ul.checkbox-list {
  columns: 2;
  padding: 0;
}
ul.checkbox-list > li {
  list-style: none;
}
</style>
