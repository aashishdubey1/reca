import type { LucideIcon } from 'lucide-react'

export type Rating = 'Hard' | 'Medium' | 'Easy'

export interface TopicItem {
  id: number
  title: string
  category: string
  icon?: LucideIcon
  tone?: 'orange' | 'neutral' | string
  due: string
  status: 'pending' | 'completed'
  notes?: string
  easeFactor?: number
  intervalDays?: number
}

export interface DayData {
  offset: number
  label: string
  date: string
  count: string
  topics: string[]
  side: 'left' | 'center' | 'right'
}
