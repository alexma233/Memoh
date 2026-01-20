import { defineStore } from 'pinia'
import { reactive,watch } from 'vue'
import { useLocalStorage } from '@vueuse/core'
import { useRouter } from 'vue-router'


type user={
  'id': string,
  'username': string,
  'role': string,
  'displayName': string
}


export const useUserStore = defineStore('user', () => {
  const userInfo = reactive<user>({
    'id': '',
    'username': '',
    'role': '',
    'displayName': ''
  })

  const localToken=useLocalStorage('token','')
  

  const login = (userData: user,token:string) => {
    localToken.value=token
    for (const key of Object.keys(userData) as (keyof user)[]) {
      userInfo[key] = userData[key]
    }
  }

  const exitLogin = () => {
    localToken.value=''
    for (const key of Object.keys(userInfo) as (keyof user)[]) {
      userInfo[key]=''
    }
  }
  const router=useRouter()
  watch(localToken, () => {
   
    if (!localToken.value) {
      exitLogin()
      router.replace({name:'Login'})
    }
  }, {
    immediate: true
  })
  return {
    userInfo,
    login,
    exitLogin
  }
}, {
  persist:true
})
