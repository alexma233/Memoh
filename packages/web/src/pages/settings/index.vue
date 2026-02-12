<template>
  <section class="h-full max-w-7xl mx-auto p-6">
    <section class="min-w-0">
      <div class="max-w-3xl mx-auto">
        <h6 class="mb-2 flex items-center">
          <FontAwesomeIcon
            :icon="['fas', 'gear']"
            class="mr-2"
          />
          {{ $t('settings.display') }}
        </h6>
        <Separator />

        <div class="mt-4 space-y-4">
          <div class="flex items-center justify-between">
            <Label>{{ $t('settings.language') }}</Label>
            <Select
              :model-value="language"
              @update:model-value="(v) => v && setLanguage(v as Locale)"
            >
              <SelectTrigger class="w-40">
                <SelectValue :placeholder="$t('settings.languagePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="zh">
                    {{ $t('settings.langZh') }}
                  </SelectItem>
                  <SelectItem value="en">
                    {{ $t('settings.langEn') }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>

          <Separator />

          <div class="flex items-center justify-between">
            <Label>{{ $t('settings.theme') }}</Label>
            <Select
              :model-value="theme"
              @update:model-value="(v) => v && setTheme(v as 'light' | 'dark')"
            >
              <SelectTrigger class="w-40">
                <SelectValue :placeholder="$t('settings.themePlaceholder')" />
              </SelectTrigger>
              <SelectContent>
                <SelectGroup>
                  <SelectItem value="light">
                    {{ $t('settings.themeLight') }}
                  </SelectItem>
                  <SelectItem value="dark">
                    {{ $t('settings.themeDark') }}
                  </SelectItem>
                </SelectGroup>
              </SelectContent>
            </Select>
          </div>
        </div>
      </div>
    </section>
  </section>
</template>

<script setup lang="ts">
import {
  Select,
  SelectTrigger,
  SelectContent,
  SelectValue,
  SelectGroup,
  SelectItem,
  Label,
  Separator,
} from '@memoh/ui'
import { storeToRefs } from 'pinia'
import { useSettingsStore } from '@/store/settings'
import type { Locale } from '@/i18n'

const settingsStore = useSettingsStore()
const { language, theme } = storeToRefs(settingsStore)
const { setLanguage, setTheme } = settingsStore
</script>
