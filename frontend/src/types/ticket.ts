export interface Ticket {
  id: number
  user_id: number
  event_id: number
  status: 'active' | 'cancelled' | 'transferred'
  purchase_price: number
  purchased_at: string
  event_title: string
  event_date: string
}

export interface PurchaseRequest {
  event_id: number
  presale_code?: string
}

export interface TransferRequest {
  ticket_id: number
  to_user_email: string
}
