import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

// Mock localStorage globally BEFORE any module import
const store: Record<string, string> = {}
const localStorageMock = {
  getItem: vi.fn((key: string) => store[key] ?? null),
  setItem: vi.fn((key: string, val: string) => { store[key] = val }),
  removeItem: vi.fn((key: string) => { delete store[key] }),
  clear: vi.fn(() => { Object.keys(store).forEach(k => delete store[k]) }),
  get length() { return Object.keys(store).length },
  key: vi.fn((i: number) => Object.keys(store)[i] ?? null),
}
vi.stubGlobal('localStorage', localStorageMock)

vi.stubGlobal('matchMedia', vi.fn().mockImplementation((query: string) => ({
  matches: false, media: query, onchange: null,
  addListener: vi.fn(), removeListener: vi.fn(),
  addEventListener: vi.fn(), removeEventListener: vi.fn(),
  dispatchEvent: vi.fn(),
})))

const { useAppStore } = await import('./app')

describe('useAppStore', () => {
  beforeEach(() => {
    Object.keys(store).forEach(k => delete store[k])
    vi.clearAllMocks()
    setActivePinia(createPinia())
  })

  it('initializes with default values', () => {
    const appStore = useAppStore()
    expect(appStore.cardsPerRow).toBe('5')
    expect(appStore.themeMode).toBe('system')
    expect(appStore.networkMode).toBe('internal')
    expect(appStore.showWeatherWidget).toBe(true)
    expect(appStore.showMemo).toBe(true)
  })

  it('setCardsPerRow updates value and persists to localStorage', () => {
    const appStore = useAppStore()
    appStore.setCardsPerRow('3')
    expect(appStore.cardsPerRow).toBe('3')
    expect(localStorage.setItem).toHaveBeenCalledWith(
      expect.stringContaining('cards-per-row'),
      '3'
    )
  })

  it('setTheme updates themeMode', () => {
    const appStore = useAppStore()
    appStore.setTheme('dark')
    expect(appStore.themeMode).toBe('dark')
    appStore.setTheme('light')
    expect(appStore.themeMode).toBe('light')
  })

  it('toggleNetwork switches between internal and external', () => {
    const appStore = useAppStore()
    expect(appStore.networkMode).toBe('internal')
    appStore.toggleNetwork()
    expect(appStore.networkMode).toBe('external')
    appStore.toggleNetwork()
    expect(appStore.networkMode).toBe('internal')
  })
})
