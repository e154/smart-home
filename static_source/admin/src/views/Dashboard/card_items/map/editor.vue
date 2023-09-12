<script setup lang="ts">
import {computed, PropType, ref, unref, watch} from "vue";
import {Card, CardItem, Core, Tab} from "@/views/Dashboard/core";
import {ElDivider, ElSwitch, ElCollapse, ElFormItem, ElInputNumber, ElCol, ElRow, ElButton, ElInput,
  ElPopconfirm, ElForm, ElCard, ElCollapseItem, ElTag} from 'element-plus'
import {ApiEntity, ApiImage} from "@/api/stub";
import CommonEditor from "@/views/Dashboard/card_items/common/editor.vue";
import EntitySearch from "@/views/Entities/components/EntitySearch.vue";
import ImageSearch from "@/views/Images/components/ImageSearch.vue";
import {useI18n} from "@/hooks/web/useI18n";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Nullable<Core>>,
    default: () => null
  },
  item: {
    type: Object as PropType<Nullable<CardItem>>,
    default: () => null
  },
})

const currentCore = computed(() => props.core as Core)

const currentItem = computed(() => props.item as CardItem)

// ---------------------------------
// markers
// ---------------------------------

const onSelectImageForAction = (index: number, image: ApiImage) => {
  // console.log('select image', index, image);
  if (!currentItem.value.payload.map.markers[index]) {
    return;
  }

  currentItem.value.payload.map.markers[index].image = image as ApiImage || undefined;
}

const changedForActionButton = async (entity: ApiEntity, index: number) => {
  if (!currentItem.value.payload.map?.markers[index]) {
    return
  }

  if (entity?.id) {
    currentItem.value.payload.map.markers[index].entity = await currentCore.value.fetchEntity(entity.id);
    currentItem.value.payload.map.markers[index].entityId = entity.id;
  } else {
    currentItem.value.payload.map.markers[index].entity = undefined;
    currentItem.value.payload.map.markers[index].entityId = '';
  }
}

const addMarker = () => {
  if (!currentItem.value.payload?.map?.markers) {
    currentItem.value.payload.map.markers = [];
  }
  currentItem.value.payload.map.markers.push({
    image: null,
    entityId: null,
    attribute: '',
    opacity: 0.9,
    scale: 0.5,
    value: [0,0],
  });
  currentItem.value.update()
}

const removeMarker = (index: number) => {
  if (!currentItem.value.payload?.map?.markers) {
    return;
  }

  currentItem.value.payload.map.markers.splice(index, 1);
  currentItem.value.update()
}

const updateCenter = (index: number) => {
  currentItem.value.payload.map.indexMarkerCenter = index
}

</script>

<template>
  <div class="mb-20px">

    <CommonEditor :item="item" :core="core"/>

    <ElDivider content-position="left">{{$t('dashboard.editor.mapOptions')}}</ElDivider>

    <ElFormItem :label="$t('dashboard.editor.staticPosition')" prop="round">
      <ElSwitch v-model="currentItem.payload.map.staticCenter"/>
    </ElFormItem>

    <!-- marker options -->
    <ElDivider content-position="left">{{$t('dashboard.editor.markers') }}</ElDivider>

    <ElRow>
      <ElCol>
        <div style="padding-bottom: 20px">
          <ElButton type="default" @click.prevent.stop="addMarker()">
            <Icon icon="ep:plus" class="mr-5px"/>
            {{ $t('dashboard.editor.addMarker') }}
          </ElButton>
        </div>

        <!-- props -->
        <ElCollapse>
          <ElCollapseItem
              v-for="(prop, index) in item.payload.map.markers"
              :name="index"
              :key="index"
          >

            <template #title>
              {{ prop.entityId }} - {{ prop.attribute }}&nbsp;&nbsp;
              <ElTag v-if="item.payload.map.indexMarkerCenter === index">
                {{ $t('dashboard.editor.tracked') }}
              </ElTag>
            </template>

            <ElCard shadow="never" class="item-card-editor">

              <ElForm
                  label-position="top"
                  :model="prop"
                  style="width: 100%"
                  ref="cardItemForm">

                <ElRow :gutter="24">
                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
                      <EntitySearch v-model="prop.entity" @change="changedForActionButton($event, index)"/>
                    </ElFormItem>
                  </ElCol>

                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.value')" prop="value">
                      <ElInput size="small" v-model="prop.attribute"/>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElRow :gutter="24">
                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.opacity')" prop="entity">
                      <ElInputNumber v-model="prop.opacity" :show-tooltip="false" :min="0" :max="1" :step="0.01" style="width: 100%"/>
                    </ElFormItem>
                  </ElCol>

                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.scale')" prop="value">
                      <ElInputNumber v-model="prop.scale" :show-tooltip="false" :min="0" :max="1" :step="0.01" style="width: 100%"/>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElFormItem :label="$t('dashboard.editor.image')" prop="image">
                  <ImageSearch v-model="prop.image" @change="onSelectImageForAction(index, ...arguments)"/>
                </ElFormItem>

                <ElRow>
                  <ElCol>
                    <div style="padding-bottom: 20px">
                      <div style="text-align: right;">
                        <ElButton
                            class="mr-10px" plain
                            @click="updateCenter(index)"
                            :disabled="currentItem.payload.map.indexMarkerCenter === index"
                        >
                          {{$t('dashboard.editor.followTheMarker')}}
                        </ElButton>
                        <ElPopconfirm
                            :confirm-button-text="$t('main.ok')"
                            :cancel-button-text="$t('main.no')"
                            width="250"
                            style="margin-left: 10px;"
                            :title="$t('main.are_you_sure_to_do_want_this?')"
                            @confirm="removeMarker(index)"
                        >
                          <template #reference>
                            <ElButton class="mr-10px" type="danger" plain>
                              <Icon icon="ep:delete" class="mr-5px"/>
                              {{ t('main.remove') }}
                            </ElButton>
                          </template>
                        </ElPopconfirm>
                      </div>
                    </div>
                  </ElCol>
                </ElRow>

              </ElForm>

            </ElCard>

          </ElCollapseItem>
        </ElCollapse>
        <!-- /props -->

      </ElCol>
    </ElRow>
    <!-- /button options -->

  </div>
</template>

<style lang="less" >

</style>
