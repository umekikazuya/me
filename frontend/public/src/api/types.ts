export interface ProblemInvalidParam {
  name?: string
  reason?: string
}

export interface LegacyProblemDetailField {
  field?: string
  message?: string
}

export interface ProblemDetail {
  type?: string
  title?: string
  status?: number
  detail?: string
  instance?: string
  invalidParams?: ProblemInvalidParam[]
  code?: string
  message?: string
  details?: LegacyProblemDetailField[]
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

export interface AdminLoginInput {
  emailAddress: string
  password: string
}

export interface ChangeEmailInput {
  token: string
  newEmailAddress: string
}

export interface MeSkillGroup {
  category: string
  items: string[]
  sortOrder: number
}

export interface MeCertification {
  name: string
  issuer: string
  year: number
  month?: number
}

export interface MeExperience {
  company: string
  url: string
  startYear: number
  endYear?: number
}

export interface MeLink {
  platform: string
  url: string
}

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

const normalizeSkillGroup = (value: unknown): MeSkillGroup => {
  const item = isRecord(value) ? value : {}
  return {
    category: asString(item.category),
    items: asStringArray(item.items),
    sortOrder: asNumber(item.sortOrder) ?? 0,
  }
}

const normalizeCertification = (value: unknown): MeCertification => {
  const item = isRecord(value) ? value : {}
  return {
    name: asString(item.name),
    issuer: asString(item.issuer),
    year: asNumber(item.year) ?? new Date().getFullYear(),
    month: asNumber(item.month),
  }
}

const normalizeExperience = (value: unknown): MeExperience => {
  const item = isRecord(value) ? value : {}
  return {
    company: asString(item.company),
    url: asString(item.url),
    startYear: asNumber(item.startYear) ?? new Date().getFullYear(),
    endYear: asNumber(item.endYear),
  }
}

const normalizeLink = (value: unknown): MeLink => {
  const item = isRecord(value) ? value : {}
  return {
    platform: asString(item.platform),
    url: asString(item.url),
  }
}

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

export const cloneMeProfile = (profile: MeProfile): MeProfile =>
  structuredClone(profile)

export const normalizeMeResponse = (payload: unknown): MeProfile => {
  const record = isRecord(payload) ? payload : {}

  return {
    displayName: asString(record.displayName),
    displayJa: asString(record.displayJa),
    role: asString(record.role),
    location: asString(record.location),
    skills: Array.isArray(record.skills)
      ? record.skills.map(normalizeSkillGroup)
      : [],
    certifications: Array.isArray(record.certifications)
      ? record.certifications.map(normalizeCertification)
      : [],
    experiences: Array.isArray(record.experiences)
      ? record.experiences.map(normalizeExperience)
      : [],
    links: Array.isArray(record.links) ? record.links.map(normalizeLink) : [],
    likes: asStringArray(record.likes),
    updatedAt: asString(record.updatedAt) || undefined,
  }
}

const trimOptional = (value: string) => {
  const trimmed = value.trim()
  return trimmed === '' ? undefined : trimmed
}

export const toMeRequest = (profile: MeProfile) => ({
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

const describeInvalidParam = (param: ProblemInvalidParam) => {
  if (param.name && param.reason) return `${param.name}: ${param.reason}`
  return param.reason || param.name
}

const describeLegacyDetail = (detail: LegacyProblemDetailField) => {
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
    ...(problem?.details?.map(describeLegacyDetail).filter(Boolean) ?? []),
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
