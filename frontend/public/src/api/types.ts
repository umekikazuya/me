import type { components } from '@me/types'

/**
 * RFC 9457 の ProblemDetails を緩めた受信用の型。
 * 生成型の必須性はワイヤの契約であり、壊れたレスポンスも受け取りうるので
 * クライアント側では optional に緩める。
 */
export type ProblemDetail = Partial<
  Omit<components['schemas']['ProblemDetails'], 'invalidParams'>
> & {
  invalidParams?: Partial<components['schemas']['InvalidParam']>[]
}

export class ApiError extends Error {
  readonly status: number
  readonly problem?: ProblemDetail

  constructor(message: string, status: number, problem?: ProblemDetail) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.problem = problem
  }
}

const describeInvalidParam = (
  param: Partial<components['schemas']['InvalidParam']>,
) => {
  if (param.name && param.reason) return `${param.name}: ${param.reason}`
  return param.reason || param.name
}

export const describeProblemDetail = (
  problem?: ProblemDetail,
  fallbackStatus?: number,
) => {
  const fieldMessages =
    problem?.invalidParams?.map(describeInvalidParam).filter(Boolean) ?? []

  return (
    problem?.detail ||
    fieldMessages.join('\n') ||
    problem?.title ||
    (fallbackStatus
      ? `API request failed with status ${fallbackStatus}`
      : 'API request failed')
  )
}

export const describeApiError = (error: unknown) => {
  if (error instanceof ApiError) {
    return describeProblemDetail(error.problem, error.status)
  }

  if (error instanceof Error) return error.message
  return '予期しないエラーが発生しました。'
}
