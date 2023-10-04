<script setup lang="ts">
import { ElDropdown, ElDropdownMenu, ElDropdownItem, ElMessageBox, ElImage } from 'element-plus'
import { useI18n } from '@/hooks/web/useI18n'
import { useCache } from '@/hooks/web/useCache'
import { resetRouter } from '@/router'
import { useRouter } from 'vue-router'
import { useDesign } from '@/hooks/web/useDesign'
import { useTagsViewStore } from '@/store/modules/tagsView'
import {useAppStore} from "@/store/modules/app";

const tagsViewStore = useTagsViewStore()
const appStore = useAppStore()
const { getPrefixCls } = useDesign()
const prefixCls = getPrefixCls('user-info')
const { t } = useI18n()
const { wsCache } = useCache()
const { replace, push } = useRouter()

const loginOut = () => {
  ElMessageBox.confirm(t('common.loginOutMessage'), t('common.reminder'), {
    confirmButtonText: t('common.ok'),
    cancelButtonText: t('common.cancel'),
    type: 'warning'
  })
    .then(async () => {

      appStore.RemoveToken()
      //wsCache.clear()

      tagsViewStore.delAllViews()
      resetRouter() // 重置静态路由表
      replace('/login')
      location.reload() // To prevent bugs from vue-router
    })
    .catch(() => {})
}

const toDocument = () => {
  const user = appStore.getUser
  push(`/etc/users/edit/${user.id}`)
}

const getAvatar = (): string => {
  return appStore.getAvatar || 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif'
}

const getUserName = (): string => {
  return appStore.getUser.nickname || 'unknown'
}

</script>

<template>
  <ElDropdown :class="prefixCls" trigger="click">
    <div class="flex items-center">
      <img
        :src="getAvatar()"
        alt="avatar"
        class="w-[calc(var(--logo-height)-10px)] rounded-[20%]"
        style="aspect-ratio: 1 / 1;"
      />
      <span class="<lg:hidden text-14px pl-[5px] text-[var(--top-header-text-color)]">{{getUserName()}}</span>
    </div>
    <template #dropdown>
      <ElDropdownMenu>
        <ElDropdownItem>
          <div @click="toDocument">{{ t('common.userProfile') }}</div>
        </ElDropdownItem>
        <ElDropdownItem divided>
          <div @click="loginOut">{{ t('common.loginOut') }}</div>
        </ElDropdownItem>
      </ElDropdownMenu>
    </template>
  </ElDropdown>
</template>
