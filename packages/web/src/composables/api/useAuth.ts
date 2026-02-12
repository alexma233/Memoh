import { postAuthLogin } from '@memoh/sdk'
import type { HandlersLoginRequest, HandlersLoginResponse } from '@memoh/sdk'

export type LoginRequest = HandlersLoginRequest
export type LoginResponse = HandlersLoginResponse

export async function login(data: LoginRequest): Promise<LoginResponse> {
  const { data: res } = await postAuthLogin({
    body: data,
    throwOnError: true,
  })
  return res
}
