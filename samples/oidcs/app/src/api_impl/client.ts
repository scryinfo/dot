import { createConnectTransport } from '@connectrpc/connect-web'
import { createClient /*, type Interceptor*/ } from '@connectrpc/connect'
import { config } from '@/config'
import { AuthService } from 'oidc_ts/src/index'

// const headerInterceptor: Interceptor = (next) => async (req) => {
//   req.header.set('Content-Type', 'application/proto')
//   return await next(req)
// }
const transport = createConnectTransport({
  baseUrl: config.API_BASE,
  useBinaryFormat: true,
  // defaultTimeoutMs: 10000,
  // fetch: (input, init) => {
  //   return fetch(input, {
  //     ...init,
  //     credentials: 'include',
  //     headers: {
  //       'Content-Type': 'application/proto',
  //     },
  //   })
  // },
  // interceptors: [headerInterceptor],
})

export const succussCode = 0

export const authService = createClient(AuthService, transport)
