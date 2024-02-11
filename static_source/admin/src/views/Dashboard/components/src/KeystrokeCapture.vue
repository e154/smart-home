<script setup lang="ts">

import {computed, onMounted, PropType, ref} from "vue";
import {Card, Core} from "@/views/Dashboard/core/core";
import {ElButton, ElCol, ElCollapse, ElCollapseItem, ElDivider, ElForm, ElPopconfirm, ElRow, ElTag, ElCard} from "element-plus";
import {useI18n} from "@/hooks/web/useI18n";
import {useEventBus} from "@/hooks/event/useEventBus";
import {EntitiesAction, EntitiesActionOptions} from "@/components/EntitiesAction";

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

onMounted(() => {
  useEventBus({
    name: 'keydown',
    callback: (val) => {
      //console.debug(val)
      if (!currentCard.value?.keysCapture) {
        return;
      }
      if (activeItemIdx.value > -1) {
        currentCard.value.keysCapture[activeItemIdx.value].keys.set(val.keyCode, val.key)
        activeItemIdx.value = -1
        return
      }
    }
  })
})

const addAction = () => {

  if (!currentCard.value?.keysCapture) {
    currentCard.value.keysCapture = []
  }

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

const changedForActionButton = async (options: EntitiesActionOptions, index: number) => {
  currentCard.value.keysCapture[index].entityId = options.entityId
  currentCard.value.keysCapture[index].entity = options.entity
  currentCard.value.keysCapture[index].action = options.action
  currentCard.value.keysCapture[index].tags = options.tags
  currentCard.value.keysCapture[index].areaId = options.areaId
}

</script>

<template>
  <ElForm
      label-position="top"
      style="width: 100%"
  >
    <ElRow>
      <ElCol>
        <ElButton class="w-[100%]" @click.prevent.stop="addAction()">
          <Icon icon="ep:plus" class="mr-5px"/>
          {{ $t('dashboard.editor.addAction') }}
        </ElButton>
        <!-- props -->
        <ElCollapse v-if="currentCard && currentCard.keysCapture">
          <ElCollapseItem
              v-for="(prop, index) in currentCard.keysCapture"
              :name="index"
              :key="index"
          >

            <template #title>
              <span>ID: {{ index }}</span>

              <span v-if="prop.entityId">
                &nbsp;
              <ElTag class="mx-1" type="info">{{ prop.entityId }}</ElTag>
                </span>

              <span v-if="prop.action">
                &nbsp;
              <ElTag class="mx-1" type="danger">{{ prop.action }}</ElTag>
                </span>
            </template>

            <ElCard>
            <ElForm
                label-position="top"
                :model="prop"
                style="width: 100%">

              <ElRow class="mt-10px mb-10px">
                <ElCol>
                  <el-button class="w-[100%]" size="small" @click="saveNewButton(index)">
                    <Icon icon="ep:plus" class="mr-5px"/>
                    {{ t('dashboard.editor.addNewButton') }}
                  </el-button>
                </ElCol>
              </ElRow>


              <ElRow class="mb-10px">
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

              <ElRow class="mb-10px">
                <ElCol>
                  <ElDivider content-position="left">{{ $t('dashboard.editor.actionOptions') }}</ElDivider>
                </ElCol>
              </ElRow>

              <ElRow class="mb-10px">
                <ElCol>
                  <EntitiesAction :options="prop" @change="changedForActionButton($event, index)"/>
                </ElCol>
              </ElRow>

              <ElRow class="mb-10px">
                <ElCol>

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
