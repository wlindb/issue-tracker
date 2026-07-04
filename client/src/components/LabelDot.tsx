import { cn } from '@/lib/utils'
import { getLabelColor } from '@/lib/labelColor'

interface LabelDotProps {
  labelId: string
  className?: string
}

export function LabelDot({ labelId, className }: LabelDotProps) {
  return (
    <span
      className={cn('size-2 rounded-full shrink-0', className)}
      style={{ backgroundColor: getLabelColor(labelId) }}
    />
  )
}
