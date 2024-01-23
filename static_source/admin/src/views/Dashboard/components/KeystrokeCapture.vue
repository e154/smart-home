<script setup lang="ts">

import {computed, PropType, ref} from "vue";
import {Card, Core} from "@/views/Dashboard/core";
import {useEmitt} from "@/hooks/web/useEmitt";
import {
  ElButton,
  ElCard,
  ElCol,
  ElCollapse,
  ElCollapseItem,
  ElForm,
  ElFormItem, ElMessage,
  ElOption,
  ElPopconfirm,
  ElRow,
  ElSelect,
  ElTag
} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {ApiEntity, ApiEntityCallActionRequest, ApiTypes} from "@/api/stub";
import EntitySearch from "@/views/Entities/components/EntitySearch.vue";
import api from "@/api/api";

const {t} = useI18n()

// ---------------------------------
// common
// ---------------------------------

const props = defineProps({
  core: {
    type: Object as PropType<Core>,
  },
  card: {
    type: Object as PropType<Nullable<Card>>,
    default: () => null
  },
})

const currentCard = computed({
  get(): Card {
    return props.card as Card
  },
  set(val: Card) {
  }
})

const currentCore = computed(() => props.core as Core)

// ---------------------------------
// component methods
// ---------------------------------

const activeItemIdx = ref(-1)

useEmitt({
  name: 'keydown',
  callback: (val) => {
    if (activeItemIdx.value > -1) {
      currentCard.value.keysCapture[activeItemIdx.value].keys.set(val.keyCode, val.key)
      activeItemIdx.value = -1
      return
    }
  }
})

const addAction = () => {
  currentCard.value.keysCapture.push({
    entity: undefined,
    entityId: currentCard.value.entityId || undefined,
    action: '',
    keys: new Map(),
  });
}

const removeAction = (index: number) => {
  if (!currentCard.value.keysCapture) {
    return;
  }

  currentCard.value.keysCapture.splice(index, 1);
}

// tags
const removeButtonHandler = (index: number, key: number) => {
  if (index > -1) {
    currentCard.value.keysCapture[index].keys.delete(key)
  }
}

const saveNewButton = (index: number) => {
  activeItemIdx.value = index
}

const changedForActionButton = async (entity: ApiEntity, index: number) => {
  if (entity?.id) {
    const _entity = await currentCore.value.fetchEntity(entity.id)
    entity.actions = _entity.actions;
    currentCard.value.keysCapture[index].entity = entity;
    currentCard.value.keysCapture[index].entityId = entity.id;
  } else {
    currentCard.value.keysCapture[index].entity = undefined;
    currentCard.value.keysCapture[index].entityId = '';
    currentCard.value.keysCapture[index].action = '';
  }
}

const getActionList = (entity?: ApiEntity) => {
  if (!entity) {
    return [];
  }
  return entity.actions;
}


</script>

<template>
  <ElForm
      label-position="top"
      style="width: 100%"
  >
    <ElRow>
      <ElCol>
        <div style="padding-bottom: 20px">
          <ElButton type="default" @click.prevent.stop="addAction()">
            <Icon icon="ep:plus" class="mr-5px"/>
            {{ $t('dashboard.editor.addAction') }}
          </ElButton>
        </div>

        <!-- props -->
        <ElCollapse>
          <ElCollapseItem
              v-for="(prop, index) in currentCard.keysCapture"
              :name="index"
              :key="index"
          >

            <template #title>
              <span>ID: {{index}}</span>

              <span v-if="prop.entityId">
                &nbsp;
              <ElTag class="mx-1" type="info">{{ prop.entityId }}</ElTag>
                </span>

              <span v-if="prop.action">
                &nbsp;
              <ElTag class="mx-1" type="danger">{{ prop.action }}</ElTag>
                </span>
            </template>

            <ElCard shadow="never" class="item-card-editor">

              <ElForm
                  label-position="top"
                  :model="prop"
                  style="width: 100%">

                <ElRow class="mb-20px">
                  <ElCol>
                    <el-button class="button-new-tag ml-1" size="small" @click="saveNewButton(index)">
                      <Icon icon="ep:plus" class="mr-5px"/>
                      {{ t('dashboard.editor.addNewButton') }}
                    </el-button>
                  </ElCol>
                </ElRow>


                <ElRow class="mb-20px">
                  <ElCol>
                    <ElTag
                        v-for="[key, value] in prop.keys"
                        :key="key"
                        class="mx-1"
                        closable
                        type=""
                        @close="removeButtonHandler(index, key)"
                    >
                      {{ value }}
                    </ElTag>
                  </ElCol>
                </ElRow>


                <ElRow :gutter="24">
                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.entity')" prop="entity">
                      <EntitySearch v-model="prop.entity" @change="changedForActionButton($event, index)"/>
                    </ElFormItem>
                  </ElCol>

                  <ElCol :span="12" :xs="12">
                    <ElFormItem :label="$t('dashboard.editor.action')" prop="action" :aria-disabled="!prop.entity">
                      <ElSelect
                          v-model="prop.action"
                          clearable
                          :placeholder="$t('dashboard.editor.selectAction')"
                          style="width: 100%"
                      >
                        <ElOption
                            v-for="item in getActionList(prop.entity)"
                            :key="item.name"
                            :label="item.name"
                            :value="item.name"/>
                      </ElSelect>
                    </ElFormItem>
                  </ElCol>
                </ElRow>

                <ElRow>
                  <ElCol>
                    <div>
                      <div style="text-align: right;">
                        <ElPopconfirm
                            :confirm-button-text="$t('main.ok')"
                            :cancel-button-text="$t('main.no')"
                            width="250"
                            style="margin-left: 10px;"
                            :title="$t('main.are_you_sure_to_do_want_this?')"
                            @confirm="removeAction(index)"
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

  </ElForm>
</template>

<style scoped lang="less">

</style>
