'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Check } from 'lucide-react'
import { CategoryIcon } from '@/components/ui/category-icon'
import { RatingPopover } from '@/components/reviews/rating-popover'
import type { Rating, TopicItem } from '@/types'

interface RevisionRowProps {
  topic: TopicItem
  done: boolean
  onOpen: () => void
  onComplete: (rating: Rating) => void
}

export function RevisionRow({ topic, done, onOpen, onComplete }: RevisionRowProps) {
  const [ratingOpen, setRatingOpen] = useState(false)

  const handleRate = (rating: Rating) => {
    setRatingOpen(false)
    onComplete(rating)
  }

  return (
    <div className={`revision-wrap ${done ? 'is-done' : ''}`}>
      <div className="revision-row">
        <button
          type="button"
          className="revision-open"
          onClick={onOpen}
          aria-label={`Open ${topic.title}`}
        >
          <CategoryIcon topic={topic} />
        </button>

        <button type="button" className="revision-copy" onClick={onOpen}>
          <strong>{topic.title}</strong>
          <small>{topic.category}</small>
        </button>

        <span className={`due-label ${topic.due === 'Due now' ? 'urgent' : ''}`}>
          <span className="clock-dot" />
          {done ? 'Completed' : topic.due}
        </span>

        <motion.button
          type="button"
          whileTap={{ scale: 0.85 }}
          className={`check-button ${done ? 'checked' : ''}`}
          onClick={() => {
            if (!done) setRatingOpen(!ratingOpen)
          }}
          aria-label={done ? `${topic.title} completed` : `Rate ${topic.title}`}
        >
          <AnimatePresence>
            {done && (
              <motion.span
                initial={{ scale: 0, rotate: -45 }}
                animate={{ scale: 1, rotate: 0 }}
                exit={{ scale: 0 }}
                transition={{ type: 'spring', damping: 15, stiffness: 300 }}
              >
                <Check size={14} />
              </motion.span>
            )}
          </AnimatePresence>
        </motion.button>
      </div>

      <AnimatePresence>
        {ratingOpen && <RatingPopover onRate={handleRate} />}
      </AnimatePresence>
    </div>
  )
}
