import { getUsersMe, putUsersMe, putUsersMePassword, getUsersMeIdentities } from '@memoh/sdk'
import type { AccountsAccount, AccountsUpdateProfileRequest, AccountsUpdatePasswordRequest, HandlersListMyIdentitiesResponse } from '@memoh/sdk'
import { fetchApi } from '@/utils/request'

// ---- Types (re-export SDK types for backward compatibility) ----

export type UserAccount = AccountsAccount
export type UpdateMyProfileRequest = AccountsUpdateProfileRequest
export type UpdateMyPasswordRequest = AccountsUpdatePasswordRequest

export interface ChannelIdentity {
  id: string
  user_id?: string
  channel: string
  channel_subject_id: string
  display_name?: string
  metadata?: Record<string, unknown>
  created_at: string
  updated_at: string
}

export type ListMyIdentitiesResponse = HandlersListMyIdentitiesResponse

export interface IssueBindCodeRequest {
  platform?: string
  ttl_seconds?: number
}

export interface IssueBindCodeResponse {
  token: string
  platform?: string
  expires_at: string
}

export async function getMyAccount(): Promise<UserAccount> {
  const { data } = await getUsersMe({ throwOnError: true })
  return data
}

export async function updateMyProfile(body: UpdateMyProfileRequest): Promise<UserAccount> {
  const { data } = await putUsersMe({ body, throwOnError: true })
  return data
}

export async function updateMyPassword(body: UpdateMyPasswordRequest): Promise<void> {
  await putUsersMePassword({ body, throwOnError: true })
}

export async function listMyIdentities(): Promise<ListMyIdentitiesResponse> {
  const { data } = await getUsersMeIdentities({ throwOnError: true })
  return data
}

// bind_codes endpoint not in SDK, keep fetchApi
export async function issueMyBindCode(data: IssueBindCodeRequest): Promise<IssueBindCodeResponse> {
  return fetchApi<IssueBindCodeResponse>('/users/me/bind_codes', {
    method: 'POST',
    body: data,
  })
}
