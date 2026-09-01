import { api } from './index'

export interface SyncNode {
  syncId: string
  type: 'folder' | 'bookmark'
  title: string
  url?: string
  parentSyncId?: string
  index: number
  createdAt: string
  updatedAt: string
  deletedAt?: string
  rev?: number
}

export interface BmsyncStatus {
  configured: boolean
  serverUrl: string
  hasSynced: boolean
  rev: number
}

export interface ChangeIntent {
  op: 'create' | 'update' | 'delete'
  syncId?: string
  type?: 'folder' | 'bookmark'
  title?: string
  url?: string
  parentSyncId?: string
  index?: number
}

export async function getStatus(): Promise<BmsyncStatus> {
  const { data } = await api.get('/bmsync/status')
  return data
}

export async function getTree(): Promise<{ rev: number; nodes: SyncNode[] }> {
  const { data } = await api.get('/bmsync/tree')
  return data
}

export async function pullTree(): Promise<{ rev: number; nodes: SyncNode[] }> {
  const { data } = await api.post('/bmsync/pull')
  return data
}

export async function pushChanges(changes: ChangeIntent[]): Promise<{ rev: number; nodes: SyncNode[] }> {
  const { data } = await api.post('/bmsync/push', { changes })
  return data
}
