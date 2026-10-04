'use client'

import { motion, AnimatePresence } from 'framer-motion'
import { ArrowRight, X } from 'lucide-react'
import type { TopicItem } from '@/types'

interface TopicDrawerProps {
  topic: TopicItem | null
  onClose: () => void
}

export function TopicDrawer({ topic, onClose }: TopicDrawerProps) {
  return (
    <AnimatePresence>
      {topic && (
        <>
          {/* Backdrop overlay */}
          <motion.div
            key="drawer-backdrop"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
            onClick={onClose}
            className="fixed inset-0 z-40 bg-black/50 backdrop-blur-xs"
          />

          {/* Drawer content */}
          <motion.aside
            key="drawer-content"
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 28, stiffness: 280 }}
            className="drawer"
          >
            <div className="drawer-head">
              <div>
                <span className="eyebrow">TOPIC DETAILS</span>
                <h2>{topic.title}</h2>
              </div>
              <button
                type="button"
                className="icon-button"
                onClick={onClose}
                aria-label="Close topic details"
              >
                <X size={18} />
              </button>
            </div>

            <div className="drawer-section">
              <span className="eyebrow">NOTES</span>
              <p>
                {topic.notes ||
                  'Review the core concepts, common patterns, and practical examples. Add notes as you learn to make future reviews more effective.'}
              </p>
            </div>

            <div className="drawer-stats">
              <div>
                <span>Current interval</span>
                <strong>{topic.intervalDays ? `${topic.intervalDays} days` : '1 day'}</strong>
              </div>
              <div>
                <span>Ease factor</span>
                <strong>{topic.easeFactor ? topic.easeFactor.toFixed(2) : '2.50'}</strong>
              </div>
              <div>
                <span>Next review</span>
                <strong>Today</strong>
              </div>
            </div>

            <div className="drawer-section">
              <span className="eyebrow">REVISION HISTORY</span>
              <div className="history-row">
                <span>Oct 2, 2026</span>
                <b className="history-badge medium-badge">Medium</b>
              </div>
              <div className="history-row">
                <span>Sep 28, 2026</span>
                <b className="history-badge easy-badge">Easy</b>
              </div>
            </div>

            <button type="button" className="drawer-action">
              Edit topic <ArrowRight size={16} />
            </button>
          </motion.aside>
        </>
      )}
    </AnimatePresence>
  )
}
