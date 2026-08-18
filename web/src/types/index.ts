export interface User {
  id: string
  username: string
  display_name: string
  avatar: string
  role: 'admin' | 'user' | 'guest'
  status: 'approved' | 'pending' | 'rejected'
  created_at: string
}

export interface PanelGroup {
  id: string
  user_id: string
  name: string
  sort_order: number
  created_at: string
  cards: Card[]
}

export interface Card {
  id: string
  group_id: string
  user_id: string
  title: string
  url: string
  url_internal?: string
  icon?: string
  icon_color?: string
  bg_color?: string
  description?: string
  open_type: 'new_tab' | 'iframe' | 'popup'
  sort_order: number
  created_at: string
}

export interface PanelData {
  groups: PanelGroup[]
}

export interface Settings {
  [key: string]: string
}
