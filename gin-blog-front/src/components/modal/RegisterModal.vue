<script setup>
import { computed, ref } from 'vue'
import UModal from '@/components/ui/UModal.vue'
import api from '@/api'
import { useAppStore } from '@/store'

const appStore = useAppStore()

const registerFlag = computed({
  get: () => appStore.registerFlag,
  set: val => appStore.setRegisterFlag(val),
})

const form = ref({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  code: '',
})

const cooldown = ref(0)
let cooldownTimer = null

function startCooldown() {
  cooldown.value = 60
  cooldownTimer = setInterval(() => {
    cooldown.value--
    if (cooldown.value <= 0) {
      clearInterval(cooldownTimer)
      cooldownTimer = null
    }
  }, 1000)
}

async function handleSendCode() {
  const email = form.value.email.trim()
  if (!email) {
    window.$message?.warning('请先输入邮箱')
    return
  }
  if (/\s/.test(form.value.email)) {
    window.$message?.warning('不可以有空格字符')
    return
  }
  const emailReg = /^([a-zA-Z0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!emailReg.test(email)) {
    window.$message?.warning('请输入正确的邮箱')
    return
  }

  try {
    await api.sendCode({ email })
    window.$message?.success('验证码已发送，请查收邮件')
    startCooldown()
  }
  catch {}
}

async function handleRegister() {
  const { username, password, confirmPassword, email, code } = form.value

  // 空格校验
  for (const val of [username, password, confirmPassword, email, code]) {
    if (/\s/.test(val)) {
      window.$message?.warning('不可以有空格字符')
      return
    }
  }

  if (!username || !password || !confirmPassword || !email || !code) {
    window.$message?.warning('请填写所有字段')
    return
  }

  if (password !== confirmPassword) {
    window.$message?.warning('两次密码必须相同')
    return
  }

  const emailReg = /^([a-zA-Z0-9])(\w|\-)+@[a-zA-Z0-9]+\.([a-zA-Z]{2,4})$/
  if (!emailReg.test(email.trim())) {
    window.$message?.warning('请输入正确的邮箱')
    return
  }

  try {
    await api.register({
      username,
      password,
      confirm_password: confirmPassword,
      email: email.trim().toLowerCase(),
      code,
    })
    window.$message?.success('注册成功！请登录')
    form.value = { username: '', password: '', confirmPassword: '', email: '', code: '' }
    registerFlag.value = false
    appStore.setLoginFlag(true)
  }
  catch {}
}

function openLogin() {
  appStore.setRegisterFlag(false)
  appStore.setLoginFlag(true)
}
</script>

<template>
  <UModal v-model="registerFlag" :width="480">
    <div class="mx-2 my-1">
      <div class="mb-4 text-xl font-bold">
        注册
      </div>
      <div class="my-5 space-y-4">
        <!-- 用户名 -->
        <div class="flex items-center">
          <span class="mr-4 inline-block w-20 shrink-0 text-right">用户名</span>
          <input
            v-model="form.username"
            placeholder="请输入用户名"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <!-- 密码 -->
        <div class="flex items-center">
          <span class="mr-4 inline-block w-20 shrink-0 text-right">密码</span>
          <input
            v-model="form.password"
            type="password"
            placeholder="请输入密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <!-- 确认密码 -->
        <div class="flex items-center">
          <span class="mr-4 inline-block w-20 shrink-0 text-right">确认密码</span>
          <input
            v-model="form.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <!-- 邮箱 -->
        <div class="flex items-center">
          <span class="mr-4 inline-block w-20 shrink-0 text-right">邮箱</span>
          <input
            v-model="form.email"
            placeholder="请输入邮箱地址"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
        </div>
        <!-- 验证码 -->
        <div class="flex items-center">
          <span class="mr-4 inline-block w-20 shrink-0 text-right">验证码</span>
          <input
            v-model="form.code"
            placeholder="请输入验证码"
            class="block w-full border-0 rounded-md p-2 text-gray-900 shadow-sm outline-none ring-1 ring-gray-300 ring-inset placeholder:text-gray-400 focus:ring-2 focus:ring-emerald"
          >
          <button
            :disabled="cooldown > 0"
            class="ml-2 shrink-0 rounded px-3 py-2 text-sm text-white transition-colors"
            :class="cooldown > 0 ? 'bg-gray-400 cursor-not-allowed' : 'bg-emerald hover:bg-emerald/80'"
            @click="handleSendCode"
          >
            {{ cooldown > 0 ? `${cooldown}s` : '发送验证码' }}
          </button>
        </div>
      </div>
      <div class="my-2 text-center">
        <button
          class="w-full rounded bg-red py-2 text-white hover:bg-orange"
          @click="handleRegister"
        >
          注册
        </button>
        <div class="mb-2 mt-6 text-left">
          已有账号？
          <button class="duration-300 hover:text-emerald" @click="openLogin">
            登录
          </button>
        </div>
      </div>
    </div>
  </UModal>
</template>
