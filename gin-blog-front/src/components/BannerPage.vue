<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

import AppFooter from './layout/AppFooter.vue'
import ULoading from '@/components/ui/ULoading.vue'

import { useAppStore } from '@/store'

const props = defineProps({
  label: {
    type: String,
    default: 'default',
  },
  showFooter: {
    type: Boolean,
    default: true,
  },
  card: {
    type: Boolean,
    default: false,
  },
  loading: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: () => useRoute().meta?.title,
  },
})

const { pageList } = storeToRefs(useAppStore())

const coverUrl = computed(() => {
  const page = pageList.value.find(e => e.label === props.label)
  return page ? page.cover : null
})
</script>

<template>
  <!-- 顶部图片：使用 img 标签显示完整图片，可随页面滚动 -->
  <div class="banner-fade-down relative w-full">
    <img v-if="coverUrl" :src="coverUrl" class="w-full block" alt="banner">
    <div v-else class="w-full h-[280px] lg:h-[400px]" style="background: grey" />
    <h1 class="animate-fade-in-down animate-duration-800 absolute inset-0 f-c-c text-3xl font-bold text-light lg:text-4xl">
      {{ props.title }}
    </h1>
  </div>
  <!-- 主体内容 -->
  <main class="mx-1 mb-10 flex-1">
    <ULoading :show="props.loading">
      <template v-if="props.card">
        <div class="card-view card-fade-up mx-auto mb-10 mt-6 max-w-[970px] min-h-[180px] py-8 lg:mt-8 lg:px-[55px]">
          <slot v-if="!props.loading" />
        </div>
      </template>
      <template v-else>
        <div class="card-fade-up mx-auto mt-6 max-w-[1150px] min-h-[400px] px-5 py-10">
          <slot />
        </div>
      </template>
    </ULoading>
  </main>
  <AppFooter v-if="props.showFooter" />
</template>
