<script setup lang="ts">
import {inject, PropType, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import {GeoJSON} from "ol/format"
import {Collection} from "ol";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})


const center = ref([40, 40]);
const projection = ref("EPSG:4326");
const zoom = ref(8);
const rotation = ref(0);

const drawEnable = ref(false);
const drawType = ref("Polygon");

const geojsonObject = ref({
  type: "FeatureCollection",
  crs: {
    type: "name",
    properties: {
      name: "EPSG:4326",
    },
  },
  features: [
    {"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[38.774659633636475,40.55751323699951],[41.33162021636963,40.73408842086792],[41.31917476654053,39.173643589019775],[38.68591070175171,39.102489948272705],[38.774659633636475,40.55751323699951]]]},"properties":null},
    {"type":"Feature","geometry":{"type":"Polygon","coordinates":[[[37.220458984375,40.69312572479248],[37.215416431427,39.98952865600586],[37.94778823852539,40.00025749206543],[37.94774532318115,40.736470222473145],[37.220458984375,40.69312572479248]]]},"properties":null}
  ],
})

const drawstart = (event) => {
  console.log(event);
  modifyEnabled.value = false;
};

const zones = ref([]);
const selectedFeatures = ref(new Collection());
const modifyEnabled = ref(true);

const drawend = (event) => {
  zones.value.push(event.feature);
  selectedFeatures.value.push(event.feature);

  modifyEnabled.value = true;
  drawEnable.value = false;

  // let parser = new GeoJSON();
  // console.log(event.feature)
  // let area = parser.writeFeatureObject(event.feature, {featureProjection: 'EPSG:4326'});
  // geojsonObject.value.features.push(area)
  // console.log(JSON.stringify(area))
};


const features = new GeoJSON().readFeatures(geojsonObject.value)
zones.value = features;
// selectedFeatures.value.push(features[0]);

const selectConditions = inject("ol-selectconditions");
const selectCondition = selectConditions.click;

function featureSelected(event) {
  console.log(event)
  modifyEnabled.value = false;
  if (event.selected.length > 0) {
    modifyEnabled.value = true;
  }
  selectedFeatures.value = event.target.getFeatures();
}

</script>

<template>
  <input type="checkbox" id="checkbox" v-model="drawEnable" />
  <label for="checkbox">Draw Enable</label>

  <select id="type" v-model="drawType">
    <option value="Polygon">Polygon</option>
    <option value="Circle">Circle</option>
  </select>

  <ol-map
      :loadTilesWhileAnimating="true"
      :loadTilesWhileInteracting="true"
      style="height: 400px"
  >
    <ol-view
        ref="view"
        :center="center"
        :rotation="rotation"
        :zoom="zoom"
        :projection="projection"
    />

    <ol-tile-layer>
      <ol-source-osm />
    </ol-tile-layer>

    <ol-vector-layer>
      <ol-source-vector :features="zones">
        <ol-interaction-modify
            v-if="modifyEnabled"
            :features="selectedFeatures"
        />
        <ol-interaction-draw
            v-if="drawEnable"
            :type="drawType"
            @drawend="drawend"
            @drawstart="drawstart"
        >
          <ol-style>
            <ol-style-stroke color="blue" :width="2"/>
            <ol-style-fill color="rgba(255, 255, 0, 0.4)"/>
          </ol-style>
        </ol-interaction-draw>
        <ol-interaction-snap v-if="modifyEnabled || drawEnable" />
      </ol-source-vector>

      <ol-style>
        <ol-style-stroke color="red" :width="2"/>
        <ol-style-fill color="rgba(255,255,255,0.1)"/>
        <ol-style-circle :radius="7">
          <ol-style-fill color="red"/>
        </ol-style-circle>
      </ol-style>
    </ol-vector-layer>

    <ol-interaction-select
        v-if="modifyEnabled || !drawEnable"
        @select="featureSelected"
        :condition="selectCondition"
        :features="selectedFeatures"
    >
      <ol-style>
        <ol-style-stroke color="green" :width="2"/>
        <ol-style-fill color="rgba(255,255,255,0.5)"/>
      </ol-style>
    </ol-interaction-select>

  </ol-map>

</template>

<style lang="less" >

</style>
