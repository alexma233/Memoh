<template>
  <section class="ml-auto">
    <Dialog v-model:open="open">
      <DialogTrigger as-child>
        <Button variant="default">
          {{ $t("button.add",{msg:"Model"}) }}
        </Button>
      </DialogTrigger>
      <DialogContent class="sm:max-w-106.25">
        <form @submit="addModel">
          <DialogHeader>
            <DialogTitle> {{ $t("button.add", { msg: "Model" }) }}</DialogTitle>
            <DialogDescription class="mb-4">
              使用不用厂商的大模型
            </DialogDescription>
          </DialogHeader>
          <div>
            <FormField
              v-slot="{ componentField }"
              name="modelId"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Model Name
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"
                    :placeholder="$t('prompt.enter',{msg:'Model Name'})"
                    v-bind="componentField"
                    autocomplete="modelId"
                  />
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
            <FormField
              v-slot="{ componentField }"
              name="baseUrl"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Base Url
                </FormLabel>
                <FormControl>
                  <Input
                    type="text"                
                    :placeholder="$t('prompt.enter', { msg: 'Base Url' })"
                    v-bind="componentField"
                    autocomplete="baseurl"
                  />
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
            <FormField
              v-slot="{ componentField }"
              name="apiKey"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Api Key
                </FormLabel>
                <FormControl>
                  <Input
                    :placeholder="$t('prompt.enter', { msg: 'Api Key' })"
                    autocomplete="apiKey"
                    v-bind="componentField"
                  />
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
            <FormField
              v-slot="{ componentField }"
              name="clientType"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Client Type
                </FormLabel>
                <FormControl>
                  <Select v-bind="componentField">
                    <SelectTrigger class="w-full">
                      <SelectValue :placeholder="$t('prompt.select',{msg:'Client Type'})" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="OpenAI">
                          OpenAI
                        </SelectItem>
                        <SelectItem value="Anthropic">
                          Anthropic
                        </SelectItem>
                        <SelectItem value="Google">
                          Google
                        </SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
            <FormField
              v-slot="{ componentField }"
              name="name"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Display Name
                </FormLabel>
                <FormControl>
                  <Input
                    :placeholder="$t('prompt.enter', { msg: 'Display Name' })"
                    autocomplete="name"
                    v-bind="componentField"
                  />
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
            <FormField
              v-slot="{ componentField }"
              name="type"
            >
              <FormItem>
                <FormLabel class="mb-2">
                  Role
                </FormLabel>
                <FormControl>
                  <Select v-bind="componentField">
                    <SelectTrigger class="w-full">
                      <SelectValue :placeholder="$t('prompt.select', { msg: 'Role' })" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectGroup>
                        <SelectItem value="chat">
                          Chat
                        </SelectItem>
                        <SelectItem value="embedding">
                          embedding
                        </SelectItem>
                      </SelectGroup>
                    </SelectContent>
                  </Select>
                </FormControl>
                <blockquote class="h-5">
                  <FormMessage />
                </blockquote>
              </FormItem>
            </FormField>
          </div>
          <DialogFooter class="mt-4">
            <DialogClose as-child>
              <Button variant="outline">
                Cancel
              </Button>
            </DialogClose>
            <Button type="submit">
              {{ $t("button.add", { msg: "Model" }) }}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  </section>
</template>

<script setup lang="ts">
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
  Input,
  Button,
  FormField,
  FormControl,
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
  FormItem,
  FormLabel,
  FormMessage
} from '@memoh/ui'
import { useForm } from 'vee-validate'
import { inject, watch, type Ref,ref } from 'vue'
import { toTypedSchema } from '@vee-validate/zod'
import z from 'zod'
import request from '@/utils/request'
import { useMutation, useQueryCache } from '@pinia/colada'


const formSchema = toTypedSchema(z.object({
  modelId:z.string().min(1),
  baseUrl: z.string().min(1),
  apiKey: z.string().min(1),
  clientType: z.string().min(1),
  name: z.string().min(1),
  type: z.string().min(1),
}))

const form = useForm({
  validationSchema: formSchema
})

const queryCache = useQueryCache()
type ModelInfoType= Parameters<(Parameters<typeof form.handleSubmit>)[0]>[0]
const { mutate: createModel } = useMutation({
  mutation: (modelInfo:ModelInfoType ) => request({
    url: '/model',
    data: {
      ...modelInfo,      
    },
    method: 'post'
  }),
  onSettled: () => { open.value = false; queryCache.invalidateQueries({ key: ['models'], exact: true })}
})

const { mutate: updateModel } = useMutation({
  mutation: (modelInfo: ModelInfoType) => request({
    url: `/model/${editInfo.value?.id}`,
    data: {
      ...modelInfo,
    },
    method: 'PUT'
  }),
  onSettled: () => { open.value = false; queryCache.invalidateQueries({ key: ['models'], exact: true }) }
})
const addModel = form.handleSubmit(async (modelInfo) => {
  if (editInfo.value?.id) {
    updateModel(modelInfo)
  } else {
    createModel(modelInfo)   
  }
 
})

const open = inject<Ref<boolean>>('open',ref(false))
const editInfo = inject('editModelInfo',ref<null|(ModelInfoType&{id:string})>(null))
watch(open, () => {  
  if (open.value && editInfo?.value) {
    form.setValues(editInfo.value) 
  }
}, {
  immediate:true
})
</script>