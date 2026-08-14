import type {
  ChangeEmailRequest,
  DomainErrorFieldDetail,
  InvalidParam,
  LoginRequest,
  MeCertification,
  MeExperience,
  MeLink,
  MeRequest,
  MeSkillGroup,
  ProblemDetails,
} from '@me/types'

export type {
  ChangeEmailRequest,
  InvalidParam,
  LoginRequest,
  MeCertification,
  MeExperience,
  MeLink,
  MeRequest,
  MeSkillGroup,
}

/**
 * エラーレスポンスの両形を受け取れる緩い型。
 *
 * サーバーは HTTP エラー (RFC 9457 の ProblemDetails) とドメインエラー (422 の
 * DomainErrorResponse) を shape で使い分けるが、クライアントは受信するまで
 * どちらか分からないため合成した形で受ける。壊れたレスポンスも受け取りうるので
 * 生成型の必須性はここで緩める。
 */
export type ProblemDetail = Partial<Omit<ProblemDetails, 'invalidParams'>> & {
  invalidParams?: Partial<InvalidParam>[]
  code?: string
  message?: string
  details?: Partial<DomainErrorFieldDetail>[]
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

/**
 * 管理画面が扱うプロフィール。
 *
 * ワイヤ型 `MeResponse` との差分は意図的:
 * - 省略されうる文字列を `''` で埋め、編集フォームが常に値を持てるようにする
 * - `updatedAt` は空文字を `undefined` に落とす
 *
 * 配列要素は形が同じなので生成型 (`MeSkillGroup` 等) をそのまま使う。
 */
export interface MeProfile {
  displayName: string
  displayJa: string
  role: string
  location: string
  skills: MeSkillGroup[]
  certifications: MeCertification[]
  experiences: MeExperience[]
  links: MeLink[]
  likes: string[]
  updatedAt?: string
}

const isRecord = (value: unknown): value is Record<string, unknown> =>
  typeof value === 'object' && value !== null

const asString = (value: unknown) => (typeof value === 'string' ? value : '')

const asNumber = (value: unknown) =>
  typeof value === 'number' && Number.isFinite(value) ? value : undefined

const asStringArray = (value: unknown) =>
  Array.isArray(value)
    ? value.filter((item): item is string => typeof item === 'string')
    : []

export const normalizeSkillGroup = (value: unknown): MeSkillGroup => {
  const item = isRecord(value) ? value : {}
  return {
    category: asString(item.category),
    items: asStringArray(item.items),
    sortOrder: asNumber(item.sortOrder) ?? 0,
  }
}

/**
 * `year` は仕様上 required。欠落は契約違反なので、実行時の今年などを捏造せず
 * `null` を返して呼び出し側で除外する。
 */
export const normalizeCertification = (
  value: unknown,
): MeCertification | null => {
  const item = isRecord(value) ? value : {}
  const year = asNumber(item.year)
  if (year === undefined) return null

  const issuer = asString(item.issuer)
  const month = asNumber(item.month)
  return {
    name: asString(item.name),
    year,
    ...(issuer ? { issuer } : {}),
    ...(month !== undefined ? { month } : {}),
  }
}

/** `startYear` は仕様上 required。欠落時は {@link normalizeCertification} と同じ扱い。 */
export const normalizeExperience = (value: unknown): MeExperience | null => {
  const item = isRecord(value) ? value : {}
  const startYear = asNumber(item.startYear)
  if (startYear === undefined) return null

  const url = asString(item.url)
  const endYear = asNumber(item.endYear)
  return {
    company: asString(item.company),
    startYear,
    ...(url ? { url } : {}),
    ...(endYear !== undefined ? { endYear } : {}),
  }
}

export const normalizeLink = (value: unknown): MeLink => {
  const item = isRecord(value) ? value : {}
  const label = asString(item.label)
  return {
    platform: asString(item.platform),
    url: asString(item.url),
    ...(label ? { label } : {}),
  }
}

const isNotNull = <T>(value: T | null): value is T => value !== null

/** 未知の配列を編集フォームが扱える形に整える。textarea の生 JSON にも使う。 */
export const normalizeSkillGroups = (value: unknown): MeSkillGroup[] =>
  Array.isArray(value) ? value.map(normalizeSkillGroup) : []

export const normalizeCertifications = (value: unknown): MeCertification[] =>
  Array.isArray(value)
    ? value.map(normalizeCertification).filter(isNotNull)
    : []

export const normalizeExperiences = (value: unknown): MeExperience[] =>
  Array.isArray(value) ? value.map(normalizeExperience).filter(isNotNull) : []

export const normalizeLinks = (value: unknown): MeLink[] =>
  Array.isArray(value) ? value.map(normalizeLink) : []

export const createEmptyMeProfile = (): MeProfile => ({
  displayName: '',
  displayJa: '',
  role: '',
  location: '',
  skills: [],
  certifications: [],
  experiences: [],
  links: [],
  likes: [],
})

export const normalizeMeResponse = (payload: unknown): MeProfile => {
  const record = isRecord(payload) ? payload : {}

  return {
    displayName: asString(record.displayName),
    displayJa: asString(record.displayJa),
    role: asString(record.role),
    location: asString(record.location),
    skills: normalizeSkillGroups(record.skills),
    certifications: normalizeCertifications(record.certifications),
    experiences: normalizeExperiences(record.experiences),
    links: normalizeLinks(record.links),
    likes: asStringArray(record.likes),
    updatedAt: asString(record.updatedAt) || undefined,
  }
}

const trimOptional = (value: string | undefined) => {
  const trimmed = value?.trim() ?? ''
  return trimmed === '' ? undefined : trimmed
}

/**
 * 送信前のサニタイズ。空文字は `undefined` にして JSON から落とし、
 * 実質的に空の行は送らない。
 *
 * 返り値に `MeRequest` を注釈しているのが要点で、spec のリクエスト契約が
 * 変わるとこの関数がコンパイルエラーになる。
 */
export const toMeRequest = (profile: MeProfile): MeRequest => ({
  displayName: profile.displayName.trim(),
  displayJa: trimOptional(profile.displayJa),
  role: trimOptional(profile.role),
  location: trimOptional(profile.location),
  skills: profile.skills
    .filter((skill) => skill.category.trim() !== '')
    .map((skill) => ({
      category: skill.category.trim(),
      items: skill.items.map((item) => item.trim()).filter(Boolean),
      sortOrder: skill.sortOrder,
    })),
  certifications: profile.certifications.map((certification) => ({
    name: certification.name.trim(),
    issuer: trimOptional(certification.issuer),
    year: certification.year,
    month: certification.month,
  })),
  experiences: profile.experiences.map((experience) => ({
    company: experience.company.trim(),
    url: trimOptional(experience.url),
    startYear: experience.startYear,
    endYear: experience.endYear,
  })),
  links: profile.links
    .filter((link) => link.platform.trim() !== '' && link.url.trim() !== '')
    .map((link) => ({
      platform: link.platform.trim(),
      url: link.url.trim(),
    })),
  likes: profile.likes.map((like) => like.trim()).filter(Boolean),
})

const describeInvalidParam = (param: Partial<InvalidParam>) => {
  if (param.name && param.reason) return `${param.name}: ${param.reason}`
  return param.reason || param.name
}

const describeDomainErrorDetail = (detail: Partial<DomainErrorFieldDetail>) => {
  if (detail.field && detail.message)
    return `${detail.field}: ${detail.message}`
  return detail.message || detail.field
}

export const describeProblemDetail = (
  problem?: ProblemDetail,
  fallbackStatus?: number,
) => {
  const fieldMessages = [
    ...(problem?.invalidParams?.map(describeInvalidParam).filter(Boolean) ??
      []),
    ...(problem?.details?.map(describeDomainErrorDetail).filter(Boolean) ?? []),
  ]

  return (
    problem?.detail ||
    fieldMessages.join('\n') ||
    problem?.title ||
    problem?.message ||
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
