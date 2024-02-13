<!--
  - This file is part of the Smart Home
  - Program complex distribution https://github.com/e154/smart-home
  - Copyright (C) 2024, Filippov Alex
  -
  - This library is free software: you can redistribute it and/or
  - modify it under the terms of the GNU Lesser General Public
  - License as published by the Free Software Foundation; either
  - version 3 of the License, or (at your option) any later version.
  -
  - This library is distributed in the hope that it will be useful,
  - but WITHOUT ANY WARRANTY; without even the implied warranty of
  - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
  - Library General Public License for more details.
  -
  - You should have received a copy of the GNU Lesser General Public
  - License along with this library.  If not, see
  - <https://www.gnu.org/licenses/>.
  -->

<script setup lang="ts">
import { Waterfall } from '@/components/Waterfall'
import { ContentWrap } from '@/components/ContentWrap'
import { useI18n } from '@/hooks/web/useI18n'
import Mock from 'mockjs'
import { ref, unref } from 'vue'
import { toAnyString } from '@/utils'

const data = ref<any>([])

const getList = () => {
  const list: any = []
  for (let i = 0; i < 20; i++) {
    // 随机 100, 500 之间的整数
    const height = Mock.Random.integer(100, 500)
    const width = Mock.Random.integer(100, 500)
    list.push(
      Mock.mock({
        width,
        height,
        id: toAnyString(),
        // http更换为https
        image_uri: Mock.Random.image(`${width}x${height}`).replace('http://', 'https://')
      })
    )
  }
  data.value = [...unref(data), ...list]
  if (unref(data).length >= 60) {
    end.value = true
  }
}
getList()

const { t } = useI18n()

const loading = ref(false)

const end = ref(false)

const loadMore = () => {
  loading.value = true
  setTimeout(() => {
    getList()
    loading.value = false
  }, 1000)
}
</script>

<template>
  <ContentWrap :title="t('router.waterfall')">
    <Waterfall
      :data="data"
      :loading="loading"
      :end="end"
      :props="{
        src: 'image_uri',
        height: 'height'
      }"
      @load-more="loadMore"
    />
  </ContentWrap>
</template>
