import { Code2, type LucideIcon } from 'lucide-react'

interface CategoryIconProps {
  topic?: {
    icon?: LucideIcon
    tone?: string
  }
  small?: boolean
}

export function CategoryIcon({ topic, small = false }: CategoryIconProps) {
  const Icon = topic?.icon || Code2
  const tone = topic?.tone || 'neutral'

  return (
    <div className={`category-icon ${tone} ${small ? 'small' : ''}`}>
      <Icon size={small ? 15 : 18} strokeWidth={1.6} />
    </div>
  )
}
