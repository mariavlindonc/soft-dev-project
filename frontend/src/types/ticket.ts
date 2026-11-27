export interface Ticket {
  id: number
  user_id: number
  event_id: number
  status: 'active' | 'cancelled' | 'transferred'
  purchase_price: number
  purchased_at: string
  event_title: string
  event_date: string
  /** true si el usuario logueado transfirió esta entrada a otra persona */
  transferred_out?: boolean
}

export interface PurchaseRequest {
  event_id: number
  quantity: number
  presale_code?: string
}

export interface TransferRequest {
  ticket_id: number
  to_user_email: string
}
