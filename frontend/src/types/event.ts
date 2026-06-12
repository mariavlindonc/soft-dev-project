export interface Event {
  id: number
  title: string
  description: string | null
  image_url: string | null
  category: string | null
  location: string | null
  event_date: string
  duration_minutes: number
  capacity: number
  tickets_sold: number
  price: number
  status: 'active' | 'presale' | 'sold_out' | 'cancelled'
  presale_active: boolean
  presale_start_date: string | null
  general_sale_date: string | null
  created_by_id: number
  created_at: string
  updated_at: string
}

export type SalePhase = 'not_yet_open' | 'presale' | 'public' | 'no_presale'

export interface EventFilters {
  category?: string
  date_from?: string
  date_to?: string
  min_price?: number
  max_price?: number
}

export interface SaleStatus {
  phase: SalePhase
  presale_start_date: string | null
  general_sale_date: string | null
  message: string
}
