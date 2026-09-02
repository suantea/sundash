import { describe, it, expect } from 'vitest'
import zhCN from './zh-CN'
import en from './en'

// Flatten all keys from a nested messages object
function flattenKeys(obj: Record<string, any>, prefix = ''): string[] {
  const keys: string[] = []
  for (const [k, v] of Object.entries(obj)) {
    const full = prefix ? `${prefix}.${k}` : k
    if (typeof v === 'object' && v !== null && !Array.isArray(v)) {
      keys.push(...flattenKeys(v, full))
    } else {
      keys.push(full)
    }
  }
  return keys
}

describe('i18n', () => {
  const zhSections = Object.keys(zhCN)
  const enSections = Object.keys(en)

  it('zh-CN and en have same top-level sections', () => {
    expect(zhSections.sort()).toEqual(enSections.sort())
  })

  it.each(zhSections)('section "%s" has same keys in zh-CN and en', (section) => {
    const zhKeys = flattenKeys(zhCN[section] as Record<string, any>)
    const enKeys = flattenKeys(en[section] as Record<string, any>)
    const missingInEn = zhKeys.filter(k => !enKeys.includes(k))
    const missingInZh = enKeys.filter(k => !zhKeys.includes(k))
    expect(missingInEn, `missing in en: ${missingInEn.join(', ')}`).toEqual([])
    expect(missingInZh, `missing in zh-CN: ${missingInZh.join(', ')}`).toEqual([])
  })

  it('no duplicate keys within a section (zh-CN)', () => {
    for (const section of zhSections) {
      const obj = zhCN[section] as Record<string, any>
      const keys = Object.keys(obj)
      const dupes = keys.filter((k, i) => keys.indexOf(k) !== i)
      expect(dupes, `duplicates in zh.${section}: ${dupes.join(', ')}`).toEqual([])
    }
  })
})
