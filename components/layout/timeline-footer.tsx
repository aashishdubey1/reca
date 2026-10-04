'use client'

import { motion } from 'framer-motion'
import { ArrowLeft, ArrowRight, CircleHelp } from 'lucide-react'

interface TimelineFooterProps {
  onScrollLeft: () => void
  onScrollRight: () => void
  onJumpToToday: () => void
}

export function TimelineFooter({
  onScrollLeft,
  onScrollRight,
  onJumpToToday,
}: TimelineFooterProps) {
  return (
    <div className="timeline-footer">
      <motion.button
        type="button"
        whileTap={{ scale: 0.9 }}
        className="bottom-arrow"
        onClick={onScrollLeft}
        aria-label="Scroll timeline left"
      >
        <ArrowLeft size={15} />
      </motion.button>

      <motion.button
        type="button"
        whileTap={{ scale: 0.95 }}
        className="jump-button"
        onClick={onJumpToToday}
      >
        Jump to today
      </motion.button>

      <motion.button
        type="button"
        whileTap={{ scale: 0.9 }}
        className="bottom-arrow"
        onClick={onScrollRight}
        aria-label="Scroll timeline right"
      >
        <ArrowRight size={15} />
      </motion.button>

      <div className="footer-hint">
        <CircleHelp size={14} /> Press ← → to navigate
      </div>
    </div>
  )
}
