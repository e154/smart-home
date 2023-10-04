<script setup lang="ts">
import {inject, PropType, reactive, ref, watch} from "vue";
import {GeoJSON} from "ol/format"
import {Collection} from "ol";
import {Fill, Stroke, Style} from "ol/style";
import {ApiArea} from "@/api/stub";
import type { ObjectEvent } from "ol/Object";
import hereIcon from "@/assets/imgs/marker.png";

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  area: {
    type: Object as PropType<Nullable<ApiArea>>,
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

watch(
    () => props.area,
    (val?: ApiArea) => {
      console.log(val)
      if (!val) {
        drawEnabled.value = true;
        return
      }
      if (val.center) {
        center.value = [val.center.lon, val.center.lat]
      }
      if (val.zoom) {
        zoom.value = val.zoom
      }
      if (val.resolution) {
        currentView.Resolution = val.resolution
      }
      let coordinates = [[]]
      for (const point of val.polygon) {
        coordinates[0].push([point.lon, point.lat])
      }

      modifyEnabled.value = val.polygon?.length > 0
      drawEnabled.value = !val.polygon || val.polygon?.length == 0

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
            "coordinates": coordinates
          },
          "properties": null
        }],
      };
      zones.value = new GeoJSON().readFeatures(geojsonObject);
    },
)

// ---------------------------------
// save/restore
// ---------------------------------

const save = () => {
  const parser = new GeoJSON();
  const colls = parser.writeFeaturesObject(zones.value, {featureProjection: 'EPSG:4326'});
  const {features} = colls
  let polygon = []
  for (const feature of features) {
    polygon = []
    for (const point of feature.geometry.coordinates[0]) {
      polygon.push({
        lon: point[0],
        lat: point[1],
      })
    }
  }
  return {
    polygon: polygon,
    zoom: currentView.Zoom,
    center: {
      lon: currentView.Center[0],
      lat: currentView.Center[1],
    },
    resolution: currentView.Resolution,
  }
}

const Clear = () => {

}


// ---------------------------------
// edit
// ---------------------------------

const zones = ref([]);
const selectedFeatures = ref(new Collection());

const drawstart = (event) => {

};

const drawend = (event) => {
  zones.value.push(event.feature);
  selectedFeatures.value.push(event.feature);

  modifyEnabled.value = true;
  drawEnabled.value = false;
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

function zoomChanged(z) {
  currentView.Zoom = z
}

function resolutionChanged(r) {
  currentView.Resolution = r
}

function centerChanged(c) {
  currentView.Center = c
}

function rotationChanged(r) {
  currentView.Rotation = r
}

const position = ref([]);
const geoLocChange = (event: ObjectEvent) => {
  // console.log("AAAAA", event);
  position.value = event.target.getPosition();
  currentView.Center = event.target?.getPosition()
};

defineExpose({
  save,
})
</script>

<template>
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
      <ol-source-osm/>
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
            type="Polygon"
            @drawstart="drawstart"
            @drawend="drawend"
        />
        <ol-interaction-snap v-if="modifyEnabled || drawEnabled"/>
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

    <ol-geolocation :projection="projection" @change:position="geoLocChange">
      <template>
        <ol-vector-layer :zIndex="2">
          <ol-source-vector>
            <ol-feature ref="positionFeature">
              <ol-geom-point :coordinates="position"/>
              <ol-style>
                <ol-style-icon :src="hereIcon" :scale="0.02"/>
              </ol-style>
            </ol-feature>
          </ol-source-vector>
        </ol-vector-layer>
      </template>
    </ol-geolocation>
  </ol-map>

<!--  <div>-->
<!--    center : {{ currentView.Center }} zoom : {{ currentView.Zoom }} resolution :-->
<!--    {{ currentView.Resolution }} rotation : {{ currentView.Rotation }}-->
<!--  </div>-->
</template>

<style lang="less">

</style>
