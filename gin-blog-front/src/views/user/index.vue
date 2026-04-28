<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

import UploadOne from './UploadOne.vue'
import BannerPage from '@/components/BannerPage.vue'

import { useUserStore } from '@/store'
import api from '@/api'

const userStore = useUserStore()
const router = useRouter()

const form = reactive({
  avatar: userStore.avatar,
  nickname: userStore.nickname,
  intro: userStore.intro,
  website: userStore.website,
  email: userStore.email,
})

// 密码修改表单
const passwordForm = reactive({
  email: '',
  code: '',
  password: '',
  confirmPassword: '',
})
const sendingCode = ref(false)
const countdown = ref(0)

onMounted(async () => {
  await userStore.getUserInfo()
  if (!userStore.userId) {
    router.push('/')
  }
  // 初始化邮箱
  passwordForm.email = userStore.email || ''
})

async function updateUserInfo() {
  try {
    await api.updateUser(form)
    window.$message?.success('修改成功!')
    userStore.getUserInfo()
  }
  catch (err) {
    console.error(err)
  }
}

async function sendCode() {
  if (!passwordForm.email) {
    window.$message?.error('请输入邮箱')
    return
  }
  sendingCode.value = true
  try {
    await api.sendCode({ email: passwordForm.email })
    window.$message?.success('验证码已发送')
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
        sendingCode.value = false
      }
    }, 1000)
  }
  catch (err) {
    sendingCode.value = false
    console.error(err)
  }
}

async function updatePassword() {
  if (!passwordForm.email) {
    window.$message?.error('请输入邮箱')
    return
  }
  if (!passwordForm.code) {
    window.$message?.error('请输入验证码')
    return
  }
  if (!passwordForm.password) {
    window.$message?.error('请输入新密码')
    return
  }
  if (passwordForm.password !== passwordForm.confirmPassword) {
    window.$message?.error('两次密码输入不一致')
    return
  }
  try {
    await api.updatePassword({
      email: passwordForm.email,
      code: passwordForm.code,
      password: passwordForm.password,
    })
    window.$message?.success('密码修改成功!')
    passwordForm.code = ''
    passwordForm.password = ''
    passwordForm.confirmPassword = ''
  }
  catch (err) {
    console.error(err)
  }
}
</script>

<template>
  <BannerPage label="user" title="个人中心" card>
    <!-- 基本信息 -->
    <p class="mb-6 text-xl font-bold">
      基本信息
    </p>
    <div class="grid grid-cols-12 gap-4">
      <div class="col-span-4 f-c-c">
        <UploadOne v-model:preview="form.avatar" />
      </div>
      <div class="col-span-8 lg:col-span-7">
        <div class="my-6 space-y-3">
          <div
            v-for="item of [
              { label: '昵称', key: 'nickname' },
              { label: '个人网站', key: 'website' },
              { label: '简介', key: 'intro' },
              { label: '邮箱', key: 'email' },
            ]" :key="item.label"
          >
            <div class="mb-2">
              {{ item.label }}
            </div>
            <input
              v-model="form[item.key]" required :placeholder="`请输入${item.label}`"
              class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
            >
          </div>
        </div>
        <button class="the-button mt-2" @click="updateUserInfo">
          修改
        </button>
      </div>
      <div class="col-span-0 lg:col-span-1" />
    </div>

    <!-- 修改密码 -->
    <p class="mb-6 mt-10 text-xl font-bold">
      修改密码
    </p>
    <div class="space-y-4">
      <div>
        <div class="mb-2">邮箱</div>
        <input
          v-model="passwordForm.email" required placeholder="请输入邮箱"
          class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
        >
      </div>
      <div>
        <div class="mb-2">验证码</div>
        <div class="flex gap-2">
          <input
            v-model="passwordForm.code" required placeholder="请输入验证码"
            class="flex-1 border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
          <button
            class="the-button shrink-0" :disabled="sendingCode"
            @click="sendCode"
          >
            {{ sendingCode ? `${countdown}s` : '发送验证码' }}
          </button>
        </div>
      </div>
      <div>
        <div class="mb-2">新密码</div>
        <input
          v-model="passwordForm.password" type="password" required placeholder="请输入新密码"
          class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
        >
      </div>
      <div>
        <div class="mb-2">确认密码</div>
        <input
          v-model="passwordForm.confirmPassword" type="password" required placeholder="请再次输入新密码"
          class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
        >
      </div>
      <button class="the-button mt-2" @click="updatePassword">
        修改密码
      </button>
    </div>
  </BannerPage>
</template>
