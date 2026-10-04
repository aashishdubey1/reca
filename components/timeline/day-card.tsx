'use client'

import { motion } from 'framer-motion'
import { BookOpen, Check, Code2 } from 'lucide-react'
import { CategoryIcon } from '@/components/ui/category-icon'
import { ProgressRing } from '@/components/ui/progress-ring'
import { RevisionRow } from '@/components/reviews/revision-row'
import { NewTopicBox } from '@/components/topics/new-topic-box'
import type { DayData, Rating, TopicItem } from '@/types'

interface DayCardProps {
  day: DayData
  centerTopics: TopicItem[]
  completed: number[]
  onComplete: (id: number, rating?: Rating) => void
  onOpenTopic: (topic: TopicItem) => void
  onAddTopic?: (title: string) => void
}

export function DayCard({
  day,
  centerTopics,
  completed,
  onComplete,
  onOpenTopic,
  onAddTopic,
}: DayCardProps) {
  const isCenter = day.side === 'center'
  const isPast = day.side === 'left'

  const countParts = day.count.split('/')
  const total = Number(countParts[1]) || 8
  const value = isCenter
    ? Math.min(completed.length, centerTopics.length)
    : isPast
      ? Number(countParts[0]) || 0
      : 0

  return (
    <motion.article
      whileHover={{ y: -4 }}
      transition={{ type: 'spring', stiffness: 350, damping: 28 }}
      className={`day-card ${day.side}`}
    >
      <div className="day-heading">
        <div>
          <h2>{day.label}</h2>
          <p>{day.date}</p>
        </div>
        {isCenter ? (
          <span className="day-pill">Day 0 · Review</span>
        ) : (
          <span className="interval-pill">
            {isPast ? 'Completed' : `Day ${day.offset} · daily view`}
          </span>
        )}
      </div>

      <div className="ring-wrap">
        <ProgressRing value={value} total={isCenter ? centerTopics.length : total} />
      </div>

      {isCenter ? (
        <>
          <div className="section-title">
            <h3>Revisions</h3>
            <span>{Math.max(0, centerTopics.length - completed.length)} remaining</span>
          </div>

          <div className="revision-list">
            {centerTopics.map((topic) => (
              <RevisionRow
                key={topic.id}
                topic={topic}
                done={completed.includes(topic.id)}
                onOpen={() => onOpenTopic(topic)}
                onComplete={(rating) => onComplete(topic.id, rating)}
              />
            ))}
          </div>

          <NewTopicBox onAddTopic={onAddTopic} />
        </>
      ) : (
        <div className="mini-topic-list">
          {day.topics.map((name, index) => (
            <div className="mini-topic" key={`${name}-${index}`}>
              <CategoryIcon
                topic={{
                  tone: isPast ? 'orange' : 'neutral',
                  icon: index % 2 ? BookOpen : Code2,
                }}
                small
              />
              <span>{name}</span>
              <i className={isPast && index < 4 ? 'done-dot' : ''}>
                {isPast && index < 4 ? <Check size={11} /> : ''}
              </i>
            </div>
          ))}
        </div>
      )}

      {isPast && (
        <div className="completed-note">
          Completed
          <br />
          <strong>{Math.abs(day.offset)} day{Math.abs(day.offset) === 1 ? '' : 's'} ago</strong>
        </div>
      )}
    </motion.article>
  )
}
