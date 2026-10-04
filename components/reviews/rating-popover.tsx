'use client'

import { motion } from 'framer-motion'
import type { Rating } from '@/types'

interface RatingPopoverProps {
  onRate: (rating: Rating) => void
}

export function RatingPopover({ onRate }: RatingPopoverProps) {
  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.88, y: 8 }}
      animate={{ opacity: 1, scale: 1, y: 0 }}
      exit={{ opacity: 0, scale: 0.92, y: 4 }}
      transition={{ duration: 0.18, ease: [0.16, 1, 0.3, 1] }}
      className="rating-popover"
      role="menu"
    >
      <span className="popover-label">How did it go?</span>
      <button type="button" onClick={() => onRate('Hard')}>
        <span className="rating-dot hard" />
        Hard
      </button>
      <button type="button" onClick={() => onRate('Medium')}>
        <span className="rating-dot medium" />
        Medium
      </button>
      <button type="button" onClick={() => onRate('Easy')}>
        <span className="rating-dot easy" />
        Easy
      </button>
    </motion.div>
  )
}
