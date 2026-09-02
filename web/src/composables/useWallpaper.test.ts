import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'

const store: Record<string, string> = {}
vi.stubGlobal('localStorage', {
  getItem: vi.fn((k: string) => store[k] ?? null),
  setItem: vi.fn((k: string, v: string) => { store[k] = v }),
  removeItem: vi.fn((k: string) => { delete store[k] }),
  clear: vi.fn(() => { Object.keys(store).forEach(k => delete store[k]) }),
  get length() { return Object.keys(store).length },
  key: vi.fn((i: number) => Object.keys(store)[i] ?? null),
})
vi.stubGlobal('matchMedia', vi.fn().mockImplementation((q: string) => ({
  matches: false, media: q, onchange: null,
  addListener: vi.fn(), removeListener: vi.fn(),
  addEventListener: vi.fn(), removeEventListener: vi.fn(),
  dispatchEvent: vi.fn(),
})))

const { useWallpaper } = await import('./useWallpaper')

describe('useWallpaper', () => {
  beforeEach(() => {
    Object.keys(store).forEach(k => delete store[k])
    setActivePinia(createPinia())
  })

  it('returns empty wallpaper style when no wallpaper set', () => {
    const { wallpaperStyle } = useWallpaper()
    expect(wallpaperStyle.value).toEqual({})
  })

  it('fetchWallpaper resolves without error when no type set', async () => {
    const { fetchWallpaper } = useWallpaper()
    await expect(fetchWallpaper()).resolves.not.toThrow()
  })

  it('wallpaperReady transitions to true after fetchWallpaper for default type', async () => {
    const { wallpaperReady, fetchWallpaper } = useWallpaper()
    expect(wallpaperReady.value).toBe(false)
    await fetchWallpaper()
    expect(wallpaperReady.value).toBe(true)
  })
})
