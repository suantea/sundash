import { api as axios } from './index'

export interface RSSFeed {
  id: string
  user_id: string
  title: string
  url: string
  description?: string
  image_url?: string
  last_fetched?: string
  update_interval: number
  created_at: string
  updated_at: string
  items?: RSSItem[]
}

export interface RSSItem {
  id: string
  feed_id: string
  title: string
  link: string
  description?: string
  pub_date?: string
  author?: string
  guid?: string
  created_at: string
  updated_at: string
}

/**
 * Get all RSS feeds for the current user
 */
export const getRSSFeeds = () => {
  return axios.get<RSSFeed[]>('rss')
}

/**
 * Add a new RSS feed
 */
export const addRSSFeed = (url: string) => {
  return axios.post<RSSFeed>('rss', { url })
}

/**
 * Update an RSS feed (e.g., change URL)
 */
export const updateRSSFeed = (feedId: string, url: string) => {
  return axios.put<RSSFeed>(`rss/${feedId}`, { url })
}

/**
 * Delete an RSS feed
 */
export const deleteRSSFeed = (feedId: string) => {
  return axios.delete(`rss/${feedId}`)
}

/**
 * Get items for a specific feed
 */
export const getRSSFeedItems = (feedId: string, limit = 10) => {
  return axios.get<RSSItem[]>(`rss/${feedId}/items`, { params: { limit } })
}