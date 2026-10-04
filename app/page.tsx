'use client'

import { useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { ArrowLeft, ArrowRight, Bell, BookOpen, Code2, Zap } from 'lucide-react'
import { DayCard } from '@/components/timeline/day-card'
import { TimelineFooter } from '@/components/layout/timeline-footer'
import { Topbar } from '@/components/layout/topbar'
import { TopicDrawer } from '@/components/topics/topic-drawer'
import type { DayData, Rating, TopicItem } from '@/types'

const initialTopics: TopicItem[] = [
  {
    id: 1,
    title: 'Go Concurrency',
    category: 'Go · Backend',
    icon: Code2,
    tone: 'orange',
    due: 'Due now',
    status: 'pending',
    intervalDays: 1,
    easeFactor: 2.5,
  },
  {
    id: 2,
    title: 'Docker Compose',
    category: 'DevOps',
    icon: Code2,
    tone: 'orange',
    due: 'Due now',
    status: 'pending',
    intervalDays: 1,
    easeFactor: 2.5,
  },
  {
    id: 3,
    title: 'Postgres Indexing',
    category: 'Database',
    icon: BookOpen,
    tone: 'neutral',
    due: 'In 2 hours',
    status: 'pending',
    intervalDays: 3,
    easeFactor: 2.5,
  },
  {
    id: 4,
    title: 'System Design Basics',
    category: 'System Design',
    icon: Zap,
    tone: 'neutral',
    due: 'In 4 hours',
    status: 'pending',
    intervalDays: 6,
    easeFactor: 2.6,
  },
]

const reviewIntervals = [1, 3, 7, 14, 30]
const extraTopics = [
  'HTTP & TCP',
  'System Design',
  'DSA',
  'Redis',
  'Operating Systems',
  'Kubernetes',
  'Distributed Systems',
  'API Design',
]

function getDayData(offset: number): DayData {
  const date = new Date()
  date.setDate(date.getDate() + offset)
  const dateLabel = date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })

  const label =
    offset === 0
      ? 'Today'
      : offset === -1
        ? 'Yesterday'
        : offset === 1
          ? 'Tomorrow'
          : offset < 0
            ? `${Math.abs(offset)} days ago`
            : `In ${offset} days`

  const scheduled = [...initialTopics.map((topic) => topic.title), ...extraTopics].filter(
    (_, index) =>
      offset >= 0 &&
      reviewIntervals.some((interval) => (offset - interval - index) % 14 === 0)
  )

  const names = scheduled.length
    ? scheduled.slice(0, 8)
    : offset < 0
      ? extraTopics.slice(0, 4)
      : ['No reviews scheduled']

  return {
    offset,
    label,
    date: dateLabel,
    count: `${offset === 0 ? 6 : offset < 0 ? 8 : 0}/${Math.max(names.length, offset === 0 ? 12 : 8)}`,
    topics: names,
    side: offset === 0 ? 'center' : offset < 0 ? 'left' : 'right',
  }
}

const initialDays = Array.from({ length: 22 }, (_, index) => getDayData(index - 7))

export default function Page() {
  const [topics, setTopics] = useState<TopicItem[]>(initialTopics)
  const [completed, setCompleted] = useState<number[]>([])
  const [drawerTopic, setDrawerTopic] = useState<TopicItem | null>(null)
  const [search, setSearch] = useState('')
  const [days, setDays] = useState<DayData[]>(initialDays)
  const [toastMessage, setToastMessage] = useState<string | null>(null)

  const timelineRef = useRef<HTMLDivElement>(null)
  const isUpdatingRef = useRef(false)

  // Drag-to-scroll refs
  const isDraggingRef = useRef(false)
  const startXRef = useRef(0)
  const scrollLeftRef = useRef(0)

  const visibleTopics = useMemo(() => {
    return topics.filter(
      (topic) =>
        topic.title.toLowerCase().includes(search.toLowerCase()) ||
        topic.category.toLowerCase().includes(search.toLowerCase())
    )
  }, [topics, search])

  const scrollTimeline = (amount: number) => {
    timelineRef.current?.scrollBy({ left: amount, behavior: 'smooth' })
  }

  const extendTimeline = () => {
    if (isUpdatingRef.current) return
    isUpdatingRef.current = true

    setDays((current) => {
      const lastOffset = current[current.length - 1]?.offset ?? 0
      const newDays = Array.from({ length: 6 }, (_, index) =>
        getDayData(lastOffset + index + 1)
      )
      return [...current, ...newDays]
    })

    setTimeout(() => {
      isUpdatingRef.current = false
    }, 200)
  }

  const prependTimeline = () => {
    if (isUpdatingRef.current) return
    isUpdatingRef.current = true

    const timeline = timelineRef.current
    if (!timeline) {
      isUpdatingRef.current = false
      return
    }

    const prevScrollWidth = timeline.scrollWidth
    const prevScrollLeft = timeline.scrollLeft

    setDays((current) => {
      const firstOffset = current[0]?.offset ?? 0
      const newDays = Array.from({ length: 6 }, (_, index) =>
        getDayData(firstOffset - 6 + index)
      )
      return [...newDays, ...current]
    })

    requestAnimationFrame(() => {
      if (timeline) {
        const widthDiff = timeline.scrollWidth - prevScrollWidth
        timeline.scrollLeft = prevScrollLeft + widthDiff
      }
      setTimeout(() => {
        isUpdatingRef.current = false
      }, 200)
    })
  }


  const jumpToToday = () => {
    const timeline = timelineRef.current
    const today = timeline?.querySelector<HTMLElement>('.day-card.center')
    if (timeline && today) {
      timeline.scrollTo({
        left: today.offsetLeft - (timeline.clientWidth - today.offsetWidth) / 2,
        behavior: 'smooth',
      })
    }
  }

  const handleAddTopic = (title: string) => {
    const newTopic: TopicItem = {
      id: Date.now(),
      title,
      category: 'General',
      tone: 'orange',
      icon: Code2,
      due: 'Due now',
      status: 'pending',
      intervalDays: 1,
      easeFactor: 2.5,
    }
    setTopics((prev) => [newTopic, ...prev])
    setToastMessage(`Added "${title}" to today's reviews`)
    setTimeout(() => setToastMessage(null), 3000)
  }

  const handleComplete = (id: number, rating?: Rating) => {
    setCompleted((current) => (current.includes(id) ? current : [...current, id]))
    const item = topics.find((t) => t.id === id)
    if (item && rating) {
      setToastMessage(`Rated ${rating}: Next review scheduled`)
      setTimeout(() => setToastMessage(null), 3000)
    }
  }

  // Scroll to Today card before first paint — no flash, no delay
  useLayoutEffect(() => {
    const timeline = timelineRef.current
    const today = timeline?.querySelector<HTMLElement>('.day-card.center')
    if (timeline && today) {
      timeline.scrollLeft = today.offsetLeft - (timeline.clientWidth - today.offsetWidth) / 2
    }
  }, [])

  // Desktop Mouse Wheel Scroll (active non-passive listener)
  useEffect(() => {
    const timeline = timelineRef.current
    if (!timeline) return

    const handleWheel = (event: WheelEvent) => {
      // Allow native vertical scroll inside the revision list
      const insideList = (event.target as HTMLElement)?.closest('.revision-list')
      if (insideList) return

      // Translate vertical wheel to horizontal carousel scroll
      if (Math.abs(event.deltaY) > Math.abs(event.deltaX)) {
        event.preventDefault()
        timeline.scrollLeft += event.deltaY * 1.15
      }
    }

    timeline.addEventListener('wheel', handleWheel, { passive: false })
    return () => {
      timeline.removeEventListener('wheel', handleWheel)
    }
  }, [])

  // Keyboard navigation
  useEffect(() => {
    const onKeyDown = (event: KeyboardEvent) => {
      if (
        event.target instanceof HTMLInputElement ||
        event.target instanceof HTMLTextAreaElement
      ) {
        return
      }

      if (event.key === 'ArrowLeft') scrollTimeline(-420)
      if (event.key === 'ArrowRight') scrollTimeline(420)
    }

    window.addEventListener('keydown', onKeyDown)
    return () => window.removeEventListener('keydown', onKeyDown)
  }, [])

  // Pointer drag-to-scroll handlers for desktop mouse
  const handlePointerDown = (e: React.PointerEvent<HTMLDivElement>) => {
    if (e.pointerType === 'mouse' && e.button === 0) {
      isDraggingRef.current = true
      startXRef.current = e.pageX - (timelineRef.current?.offsetLeft || 0)
      scrollLeftRef.current = timelineRef.current?.scrollLeft || 0
      if (timelineRef.current) {
        timelineRef.current.style.cursor = 'grabbing'
        timelineRef.current.style.userSelect = 'none'
      }
    }
  }

  const handlePointerMove = (e: React.PointerEvent<HTMLDivElement>) => {
    if (!isDraggingRef.current || !timelineRef.current) return
    e.preventDefault()
    const x = e.pageX - (timelineRef.current.offsetLeft || 0)
    const walk = (x - startXRef.current) * 1.3
    timelineRef.current.scrollLeft = scrollLeftRef.current - walk
  }

  const handlePointerUp = () => {
    isDraggingRef.current = false
    if (timelineRef.current) {
      timelineRef.current.style.cursor = 'grab'
      timelineRef.current.style.removeProperty('user-select')
    }
  }

  return (
    <main className="app-shell">
      <Topbar search={search} onSearchChange={setSearch} streak={12} />

      <section className="timeline-region">
        <motion.button
          type="button"
          whileHover={{ scale: 1.08 }}
          whileTap={{ scale: 0.94 }}
          className="nav-arrow left"
          onClick={() => scrollTimeline(-420)}
          aria-label="Previous day"
        >
          <ArrowLeft size={18} />
        </motion.button>

        <div
          ref={timelineRef}
          className="timeline"
          onPointerDown={handlePointerDown}
          onPointerMove={handlePointerMove}
          onPointerUp={handlePointerUp}
          onPointerLeave={handlePointerUp}
          onScroll={(event) => {
            const element = event.currentTarget
            if (element.scrollLeft < 120) {
              prependTimeline()
            } else if (
              element.scrollLeft + element.clientWidth >=
              element.scrollWidth - 120
            ) {
              extendTimeline()
            }
          }}
        >
          {days.map((day) => (
            <DayCard
              key={`${day.label}-${day.offset}`}
              day={day}
              centerTopics={visibleTopics}
              completed={completed}
              onComplete={handleComplete}
              onOpenTopic={setDrawerTopic}
              onAddTopic={handleAddTopic}
            />
          ))}
        </div>

        <motion.button
          type="button"
          whileHover={{ scale: 1.08 }}
          whileTap={{ scale: 0.94 }}
          className="nav-arrow right"
          onClick={() => scrollTimeline(420)}
          aria-label="Next day"
        >
          <ArrowRight size={18} />
        </motion.button>
      </section>

      <TimelineFooter
        onScrollLeft={() => scrollTimeline(-420)}
        onScrollRight={() => scrollTimeline(420)}
        onJumpToToday={jumpToToday}
      />

      <AnimatePresence>
        {toastMessage && (
          <motion.div
            key="toast-popup"
            initial={{ opacity: 0, y: 24, scale: 0.94 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 16, scale: 0.94 }}
            transition={{ type: 'spring', stiffness: 350, damping: 25 }}
            className="toast"
          >
            <Bell size={15} />
            <span>{toastMessage}</span>
          </motion.div>
        )}
      </AnimatePresence>

      <TopicDrawer topic={drawerTopic} onClose={() => setDrawerTopic(null)} />
    </main>
  )
}

