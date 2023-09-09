<script setup lang="ts">
import {computed, onMounted, PropType, ref, unref, watch, inject} from "vue";
import {Card, CardItem, Core, requestCurrentState, Tab} from "@/views/Dashboard/core";
import { GeoJSON } from "ol/format"
import { Fill, Stroke, Style } from "ol/style"
import type { ObjectEvent } from "ol/Object";
import type { View } from "ol";
import {Cache, GetTokens, RenderText} from "@/views/Dashboard/render";
import {debounce} from "lodash-es";

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
  console.log("AAAAA", event);
  position.value = event.target.getPosition();
  view.value?.setCenter(event.target?.getPosition());
};

const overviewmapcontrol = ref(true);
const fullscreencontrol = ref(true);
const mousepositioncontrol = ref(true);
const attributioncontrol = ref(true);

// ---------------------------------
// component methods
// ---------------------------------

const vectors = ref(null);
const drawedMarker = ref()
const drawType = ref("Point")

const drawstart = (event) => {
  vectors.value.source.removeFeature(drawedMarker.value);
  drawedMarker.value = event.feature;
  console.log(vectors.value.source)
}

// ---------------------------------
// component methods
// ---------------------------------

const contextMenuItems = ref([]);
const markers = ref(null);
const Feature = inject("ol-feature");
const Geom = inject("ol-geom");

contextMenuItems.value = [
  {
    text: "Center map here",
    classname: "some-style-class", // add some CSS rules
    callback: (val) => {
      view.value.setCenter(val.coordinate);
    }, // `center` is your callback function
  },
  {
    text: "Add a Marker",
    classname: "some-style-class", // you can add this icon with a CSS class
    // instead of `icon` property (see next line)
    icon: "https://upload.wikimedia.org/wikipedia/commons/thumb/8/88/Map_marker.svg/1334px-Map_marker.svg.png", // this can be relative or absolute
    callback: (val) => {
      const feature = new Feature({
        geometry: new Geom.Point(val.coordinate),
      });
      markers.value.source.addFeature(feature);
    },
  },
  "-", // this is a separator
];

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

update()

for (let index in props.item?.payload.map?.markers) {
  requestCurrentState(props.item?.payload.map?.markers[index].entityId);
}

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

<!--      <ol-overviewmap-control v-if="overviewmapcontrol">-->
<!--        <ol-tile-layer>-->
<!--          <ol-source-osm />-->
<!--        </ol-tile-layer>-->
<!--      </ol-overviewmap-control>-->

        <ol-tile-layer>
          <ol-source-osm />
        </ol-tile-layer>


<!--      <ol-tile-layer ref="jawgLayer" title="JAWG">-->
<!--        <ol-source-xyz-->
<!--            crossOrigin="anonymous"-->
<!--            url="https://c.tile.jawg.io/jawg-dark/{z}/{x}/{y}.png?access-token=87PWIbRaZAGNmYDjlYsLkeTVJpQeCfl2Y61mcHopxXqSdxXExoTLEv7dwqBwSWuJ"-->
<!--        />-->
<!--      </ol-tile-layer>-->

<!--      <ol-vector-layer>-->
<!--        <ol-source-cluster :distance="40">-->
<!--          <ol-source-vector>-->
<!--            <ol-feature v-for="index in 300" :key="index">-->
<!--              <ol-geom-point-->
<!--                  :coordinates="[-->
<!--                getRandomInRange(24, 45, 3),-->
<!--                getRandomInRange(35, 41, 3),-->
<!--              ]"-->
<!--              />-->
<!--            </ol-feature>-->
<!--          </ol-source-vector>-->
<!--        </ol-source-cluster>-->

<!--        <ol-style :overrideStyleFunction="overrideStyleFunction">-->
<!--          <ol-style-stroke color="red" :width="2"/>-->
<!--          <ol-style-fill color="rgba(255,255,255,0.1)"/>-->

<!--          <ol-style-circle :radius="10">-->
<!--            <ol-style-fill color="#3399CC"/>-->
<!--            <ol-style-stroke color="#fff" :width="1"/>-->
<!--          </ol-style-circle>-->
<!--          <ol-style-text>-->
<!--            <ol-style-fill color="#fff"/>-->
<!--          </ol-style-text>-->
<!--        </ol-style>-->
<!--      </ol-vector-layer>-->


      <ol-context-menu-control :items="contextMenuItems" />

      <ol-vector-layer>
        <ol-source-vector ref="vectors">
          <ol-interaction-draw @drawstart="drawstart" :type="drawType"/>
        </ol-source-vector>

        <ol-style>
          <ol-style-icon src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/88/Map_marker.svg/1334px-Map_marker.svg.png" :scale="0.02"/>
        </ol-style>
      </ol-vector-layer>

      <ol-vector-layer>
        <ol-source-vector ref="markers"/>
        <ol-style>
          <ol-style-icon src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/88/Map_marker.svg/1334px-Map_marker.svg.png" :scale="0.02"/>
        </ol-style>
      </ol-vector-layer>

      <ol-geolocation :projection="projection" @change:position="geoLocChange">
        <template>
          <ol-vector-layer :zIndex="2">
            <ol-source-vector>
              <ol-feature ref="positionFeature">
                <ol-geom-point :coordinates="position"/>
                <ol-style>
                  <ol-style-icon src="https://upload.wikimedia.org/wikipedia/commons/thumb/8/88/Map_marker.svg/1334px-Map_marker.svg.png" :scale="0.02"/>
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
