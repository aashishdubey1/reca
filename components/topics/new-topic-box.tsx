'use client'

import { useState } from 'react'
import { ArrowRight, Plus } from 'lucide-react'

interface NewTopicBoxProps {
  onAddTopic?: (title: string) => void
}

export function NewTopicBox({ onAddTopic }: NewTopicBoxProps) {
  const [value, setValue] = useState('')
  const [added, setAdded] = useState(false)

  const handleAdd = () => {
    const trimmed = value.trim()
    if (!trimmed) return

    if (onAddTopic) {
      onAddTopic(trimmed)
    }
    setAdded(true)
    setValue('')
    setTimeout(() => setAdded(false), 2500)
  }

  return (
    <div className="new-topic">
      <div className="new-topic-input">
        <Plus size={22} />
        <input
          value={value}
          onChange={(event) => {
            setValue(event.target.value)
            setAdded(false)
          }}
          onKeyDown={(event) => {
            if (event.key === 'Enter' && !event.nativeEvent.isComposing) {
              handleAdd()
            }
          }}
          placeholder="What do you want to learn today?"
          aria-label="New topic title"
        />
        {added && <span className="topic-added">Added</span>}
      </div>

      <div className="new-topic-actions">
        <button type="button" className="add-topic" onClick={handleAdd}>
          Add Topic <ArrowRight size={16} />
        </button>
      </div>
    </div>
  )
}
