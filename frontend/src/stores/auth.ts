import { computed, ref } from 'vue'
import { defineStore } from 'pinia'

import { ApiError, apiRequest, type ApiRequestOptions } from '@/services/api'

export interface AuthUser {
  id: number
  email: string
}

export interface TokenPair {
  access_token: string
  refresh_token: string
  token_type: string
  expires_at: string
}

export interface Profile {
  user_id: number
  first_name: string
  last_name: string
  birth_date: string
}

export interface DriverLicenseCategory {
  category_code: string
  issued_at?: string | null
  expires_at?: string | null
}

export interface DriverLicense {
  user_id: number
  license_number: string
  issued_at: string
  expires_at: string
  categories: DriverLicenseCategory[]
}

interface AuthResponse {
  user: AuthUser
  tokens: TokenPair
}

interface StoredAuthState {
  user: AuthUser | null
  tokens: TokenPair | null
  profile: Profile | null
  driverLicense: DriverLicense | null
}

interface RegisterPayload {
  email: string
  password: string
  firstName: string
  lastName: string
  birthDate: string
}

interface LoginPayload {
  email: string
  password: string
}

interface UpdateProfilePayload {
  firstName: string
  lastName: string
  birthDate: string
}

interface SaveDriverLicensePayload {
  licenseNumber: string
  issuedAt: string
  expiresAt: string
  categories: string[]
}

const STORAGE_KEY = 'three-eights-rent.auth'

let refreshPromise: Promise<boolean> | null = null

function safeParseStoredState(rawState: string | null): StoredAuthState | null {
  if (!rawState) {
    return null
  }

  try {
    return JSON.parse(rawState) as StoredAuthState
  } catch {
    return null
  }
}

function formatDateForInput(value?: string | null): string {
  return value ? value.slice(0, 10) : ''
}

export { ApiError, formatDateForInput }

export const useAuthStore = defineStore('auth', () => {
  const user = ref<AuthUser | null>(null)
  const tokens = ref<TokenPair | null>(null)
  const profile = ref<Profile | null>(null)
  const driverLicense = ref<DriverLicense | null>(null)
  const initialized = ref(false)
  const isRefreshing = ref(false)
  const isProfileLoading = ref(false)

  const isAuthenticated = computed(
    () => Boolean(user.value && tokens.value?.access_token && tokens.value?.refresh_token),
  )

  const displayName = computed(() => {
    const fullName = [profile.value?.first_name, profile.value?.last_name]
      .filter(Boolean)
      .join(' ')
      .trim()

    return fullName || user.value?.email || 'Profile'
  })

  function persistState() {
    if (!user.value || !tokens.value) {
      localStorage.removeItem(STORAGE_KEY)
      return
    }

    const snapshot: StoredAuthState = {
      user: user.value,
      tokens: tokens.value,
      profile: profile.value,
      driverLicense: driverLicense.value,
    }

    localStorage.setItem(STORAGE_KEY, JSON.stringify(snapshot))
  }

  function hydrateState() {
    const snapshot = safeParseStoredState(localStorage.getItem(STORAGE_KEY))

    if (!snapshot?.user || !snapshot.tokens) {
      clearSession()
      return
    }

    user.value = snapshot.user
    tokens.value = snapshot.tokens
    profile.value = snapshot.profile
    driverLicense.value = snapshot.driverLicense
  }

  function applyAuthResponse(response: AuthResponse) {
    user.value = response.user
    tokens.value = response.tokens
    persistState()
  }

  function clearSession() {
    user.value = null
    tokens.value = null
    profile.value = null
    driverLicense.value = null
    localStorage.removeItem(STORAGE_KEY)
  }

  async function authorizedRequest<TResponse, TBody = unknown>(
    path: string,
    options: Omit<ApiRequestOptions<TBody>, 'accessToken'> = {},
  ): Promise<TResponse> {
    if (!tokens.value?.access_token) {
      throw new ApiError('Authentication is required.', 401, 'MISSING_TOKEN')
    }

    try {
      return await apiRequest<TResponse, TBody>(path, {
        ...options,
        accessToken: tokens.value.access_token,
      })
    } catch (error) {
      if (!(error instanceof ApiError) || error.status !== 401) {
        throw error
      }

      const refreshed = await refreshTokens()
      if (!refreshed || !tokens.value?.access_token) {
        throw error
      }

      return apiRequest<TResponse, TBody>(path, {
        ...options,
        accessToken: tokens.value.access_token,
      })
    }
  }

  async function refreshTokens(): Promise<boolean> {
    if (!tokens.value?.refresh_token) {
      clearSession()
      return false
    }

    if (refreshPromise) {
      return refreshPromise
    }

    refreshPromise = (async () => {
      isRefreshing.value = true

      try {
        const nextTokens = await apiRequest<TokenPair, { refresh_token: string }>('/refresh', {
          method: 'POST',
          body: {
            refresh_token: tokens.value?.refresh_token || '',
          },
        })

        if (!user.value) {
          clearSession()
          return false
        }

        tokens.value = nextTokens
        persistState()
        return true
      } catch {
        clearSession()
        return false
      } finally {
        isRefreshing.value = false
        refreshPromise = null
      }
    })()

    return refreshPromise
  }

  async function register(payload: RegisterPayload) {
    const response = await apiRequest<
      AuthResponse,
      {
        email: string
        password: string
        first_name: string
        last_name: string
        birth_date: string
      }
    >('/register', {
      method: 'POST',
      body: {
        email: payload.email,
        password: payload.password,
        first_name: payload.firstName,
        last_name: payload.lastName,
        birth_date: payload.birthDate,
      },
    })

    applyAuthResponse(response)
    await fetchProfile()
  }

  async function login(payload: LoginPayload) {
    const response = await apiRequest<AuthResponse, { email: string; password: string }>('/login', {
      method: 'POST',
      body: {
        email: payload.email,
        password: payload.password,
      },
    })

    applyAuthResponse(response)
    await fetchProfile()
  }

  async function logout(options: { skipRequest?: boolean } = {}) {
    const refreshToken = tokens.value?.refresh_token

    try {
      if (!options.skipRequest && refreshToken) {
        await apiRequest<void, { refresh_token: string }>('/logout', {
          method: 'POST',
          body: {
            refresh_token: refreshToken,
          },
        })
      }
    } catch (error) {
      if (!(error instanceof ApiError) || error.status !== 401) {
        throw error
      }
    } finally {
      clearSession()
    }
  }

  async function fetchProfile() {
    isProfileLoading.value = true

    try {
      const profileResponse = await authorizedRequest<Profile>('/profile')
      profile.value = {
        ...profileResponse,
        birth_date: formatDateForInput(profileResponse.birth_date),
      }

      try {
        const licenseResponse = await authorizedRequest<DriverLicense>('/driver-license')
        driverLicense.value = {
          ...licenseResponse,
          issued_at: formatDateForInput(licenseResponse.issued_at),
          expires_at: formatDateForInput(licenseResponse.expires_at),
          categories: licenseResponse.categories.map((category) => ({
            ...category,
            issued_at: formatDateForInput(category.issued_at),
            expires_at: formatDateForInput(category.expires_at),
          })),
        }
      } catch (error) {
        if (error instanceof ApiError && error.status === 404) {
          driverLicense.value = null
        } else {
          throw error
        }
      }

      persistState()
    } finally {
      isProfileLoading.value = false
    }
  }

  async function updateProfile(payload: UpdateProfilePayload) {
    const response = await authorizedRequest<
      Profile,
      { first_name: string; last_name: string; birth_date: string }
    >('/profile', {
      method: 'PUT',
      body: {
        first_name: payload.firstName,
        last_name: payload.lastName,
        birth_date: payload.birthDate,
      },
    })

    profile.value = {
      ...response,
      birth_date: formatDateForInput(response.birth_date),
    }

    persistState()
  }

  async function saveDriverLicense(payload: SaveDriverLicensePayload) {
    const response = await authorizedRequest<
      DriverLicense,
      {
        license_number: string
        issued_at: string
        expires_at: string
        categories: Array<{ category_code: string; issued_at: string; expires_at: string }>
      }
    >('/driver-license', {
      method: 'POST',
      body: {
        license_number: payload.licenseNumber,
        issued_at: payload.issuedAt,
        expires_at: payload.expiresAt,
        categories: payload.categories.map((categoryCode) => ({
          category_code: categoryCode,
          issued_at: payload.issuedAt,
          expires_at: payload.expiresAt,
        })),
      },
    })

    driverLicense.value = {
      ...response,
      issued_at: formatDateForInput(response.issued_at),
      expires_at: formatDateForInput(response.expires_at),
      categories: response.categories.map((category) => ({
        ...category,
        issued_at: formatDateForInput(category.issued_at),
        expires_at: formatDateForInput(category.expires_at),
      })),
    }

    persistState()
  }

  async function initialize() {
    if (initialized.value) {
      return
    }

    hydrateState()

    if (tokens.value?.refresh_token) {
      try {
        await fetchProfile()
      } catch (error) {
        if (error instanceof ApiError && error.status === 401) {
          const refreshed = await refreshTokens()
          if (refreshed) {
            await fetchProfile()
          }
        } else {
          clearSession()
        }
      }
    }

    initialized.value = true
  }

  return {
    user,
    tokens,
    profile,
    driverLicense,
    initialized,
    isAuthenticated,
    isRefreshing,
    isProfileLoading,
    displayName,
    initialize,
    register,
    login,
    logout,
    refreshTokens,
    fetchProfile,
    updateProfile,
    saveDriverLicense,
    clearSession,
    authorizedRequest,
  }
})
