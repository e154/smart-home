<script setup lang="ts">
import {inject, PropType, reactive, ref} from "vue";
import {CardItem} from "@/views/Dashboard/core";
import {GeoJSON} from "ol/format"
import {Collection} from "ol";
import {Fill, Stroke, Style} from "ol/style";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const map = ref("");
const center = ref([-102.13121, 40.2436]);
const projection = ref("EPSG:4326");
const zoom = ref(5);
const rotation = ref(0);
const modifyEnabled = ref(false);
const drawEnabled = ref(false);
const drawType = ref("Polygon");

// ---------------------------------
// save/restore
// ---------------------------------

const fetch = async () => {
  const geojsonObject = {
    type: "FeatureCollection",
    crs: {
      type: "name",
      properties: {
        name: "EPSG:4326",
      },
    },
    features: [{
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [[[-120.57356426341968, 47.62779988411872], [-121.36271234087499, 39.45346738917534], [-113.04222608634038, 39.5500809855296], [-113.06586558332069, 47.624922032312426], [-120.57356426341968, 47.62779988411872]]]
      },
      "properties": null
    }, {
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [[[-107.95171736268756, 47.22469507039383], [-108.17433975598897, 39.07996333688665], [-99.42423133894287, 39.3753542687187], [-99.42423133894287, 47.45368984983775], [-107.95171736268756, 47.22469507039383]]]
      },
      "properties": null
    }, {
      "type": "Feature",
      "geometry": {
        "type": "Polygon",
        "coordinates": [[[-94.44143649718288, 47.70241847023914], [-94.44143649718288, 39.89254535047895], [-83.96030021864914, 40.65085930143821], [-83.96030021864914, 48.77647387651785], [-94.44143649718288, 47.70241847023914]]]
      },
      "properties": null
    }, {"type": "Feature", "geometry": {"type": "GeometryCollection", "geometries": []}, "properties": null}],
  };
  zones.value = new GeoJSON().readFeatures(geojsonObject);
}

const save = async () => {
  const parser = new GeoJSON();
  const colls = parser.writeFeaturesObject(zones.value, {featureProjection: 'EPSG:4326'});
  const {features}  = colls
  // console.log(JSON.stringify(features))
}

const deletePolygon = () => {
  selectedFeatures.value
}


const Clear = () => {

}


// ---------------------------------
// edit
// ---------------------------------

const zones = ref([]);
const selectedFeatures = ref(new Collection());

const drawstart = (event) => {
  // console.log(event);
  // modifyEnabled.value = false;
};

const drawend = (event) => {
  zones.value.push(event.feature);
  selectedFeatures.value.push(event.feature);

  modifyEnabled.value = true;
  drawEnabled.value = false;

  save()
};

function vectorStyle() {
  const style = new Style({
    stroke: new Stroke({
      color: "blue",
      width: 3,
    }),
    fill: new Fill({
      color: "rgba(0, 0, 255, 0.4)",
    }),
  });
  return style;
}

const selectConditions = inject("ol-selectconditions");
const selectCondition = selectConditions.click;

function featureSelected(event) {
  modifyEnabled.value = false;
  if (event.selected.length > 0) {
    modifyEnabled.value = true;
  }
  selectedFeatures.value = event.target.getFeatures();
  console.log(selectedFeatures.value)
}

// ---------------------------------
// view
// ---------------------------------

interface View {
  Center: Array<number>;
  Zoom: number;
  Resolution: number;
  Rotation: number;
}

const currentView = reactive<View>({
  Center: center.value,
  Zoom: zoom.value,
  Resolution: 0,
  Rotation: rotation.value,
})

function zoomChanged(z) {currentView.Zoom = z}
function resolutionChanged(r) {currentView.Resolution = r}
function centerChanged(c) {currentView.Center = c}
function rotationChanged(r) {currentView.Rotation = r}

// ---------------------------------
// etc
// ---------------------------------
fetch()

</script>

<template>
  <input type="checkbox" id="checkbox" v-model="drawEnabled" />
  <label for="checkbox">Draw Enable</label>

  <select id="type" v-model="drawType">
    <option value="Polygon">Polygon</option>
    <option value="Circle">Circle</option>
  </select>

  <ol-map
      ref="map"
      :load-tiles-while-animating="true"
      :load-tiles-while-interacting="true"
      style="height: 350px"
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

    <ol-tile-layer>
      <ol-source-osm />
    </ol-tile-layer>

    <ol-vector-layer :styles="vectorStyle">
      <ol-source-vector :features="zones">
        <ol-interaction-modify
            v-if="modifyEnabled"
            :features="selectedFeatures"
        />

        <ol-interaction-draw
            v-if="drawEnabled"
            :stopClick="true"
            :type="drawType"
            @drawstart="drawstart"
            @drawend="drawend"
        />
        <ol-interaction-snap v-if="modifyEnabled || drawEnabled" />
      </ol-source-vector>
    </ol-vector-layer>
    <ol-interaction-select
        @select="featureSelected"
        :condition="selectCondition"
        :features="selectedFeatures"
    >
      <ol-style>
        <ol-style-stroke :color="'red'" :width="2"/>
        <ol-style-fill :color="`rgba(255, 0, 0, 0.4)`"/>
      </ol-style>
    </ol-interaction-select>
  </ol-map>

  <div>
    center : {{ currentView.Center }} zoom : {{ currentView.Zoom }} resolution :
    {{ currentView.Resolution }} rotation : {{ currentView.Rotation }}
  </div>
</template>

<style lang="less" >

</style>
