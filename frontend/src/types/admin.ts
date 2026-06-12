export interface CreateEventData {
  title: string
  description?: string
  image_url?: string
  category?: string
  location?: string
  event_date: string
  duration_minutes?: number
  capacity: number
  price: number
  presale_active?: boolean
  presale_code?: string
  presale_start_date?: string
  general_sale_date?: string
}

export interface UpdateEventData {
  title?: string
  description?: string | null
  image_url?: string | null
  category?: string | null
  location?: string | null
  event_date?: string
  duration_minutes?: number | null
  capacity?: number
  price?: number
  presale_active?: boolean
  presale_code?: string | null
  presale_start_date?: string | null
  general_sale_date?: string | null
  status?: string
}

export interface EventSummary {
  event_id: number
  title: string
  capacity: number
  tickets_sold: number
  occupancy: number
}

export interface GlobalReport {
  total_events: number
  total_tickets_sold: number
  events: EventSummary[]
}

export interface BuyerInfo {
  user_id: number
  name: string
  email: string
}

export interface EventReport {
  event_id: number
  title: string
  capacity: number
  tickets_sold: number
  occupancy: number
  buyers: BuyerInfo[]
}
