'use client'

import { BarChart3, ChevronDown, Flame, Search } from 'lucide-react'

interface TopbarProps {
  search: string
  onSearchChange: (value: string) => void
  streak?: number
}

export function Topbar({ search, onSearchChange, streak = 12 }: TopbarProps) {
  return (
    <header className="topbar">
      <div className="brand">
        <span className="brand-mark">
          <span />
          <span />
          <span />
        </span>
        <span>repeats</span>
      </div>

      <div className="search-box">
        <Search size={17} />
        <input
          value={search}
          onChange={(event) => onSearchChange(event.target.value)}
          placeholder="Search topics, notes, tags..."
          aria-label="Search topics"
        />
        <kbd>⌘ K</kbd>
      </div>

      <div className="top-actions">
        <button type="button" className="top-icon" aria-label="View stats">
          <BarChart3 size={21} />
        </button>
        <span className="divider" />
        <div className="streak">
          <Flame size={21} fill="currentColor" />
          <strong>{streak}</strong>
          <span>day streak</span>
        </div>
        <span className="divider" />
        <button type="button" className="avatar" aria-label="Open profile menu">
          JD
        </button>
        <ChevronDown size={16} className="chevron" />
      </div>
    </header>
  )
}
