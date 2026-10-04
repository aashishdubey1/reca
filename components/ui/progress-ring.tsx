interface ProgressRingProps {
  value: number
  total: number
  size?: number
}

export function ProgressRing({ value, total, size = 116 }: ProgressRingProps) {
  const radius = 46
  const circumference = 2 * Math.PI * radius
  const safeTotal = total > 0 ? total : 1
  const progress = Math.min((value / safeTotal) * circumference, circumference)

  return (
    <div className="progress-ring" style={{ width: size, height: size }}>
      <svg viewBox="0 0 116 116" aria-hidden="true">
        <circle className="ring-track" cx="58" cy="58" r={radius} />
        <circle
          className="ring-value"
          cx="58"
          cy="58"
          r={radius}
          style={{
            strokeDasharray: circumference,
            strokeDashoffset: circumference - progress,
          }}
        />
      </svg>
      <div className="ring-copy">
        <strong>
          {value}/{total}
        </strong>
        <span>revisions</span>
      </div>
    </div>
  )
}
