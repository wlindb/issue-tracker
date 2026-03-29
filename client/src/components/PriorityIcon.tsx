import { AlertCircleIcon, ArrowDownIcon, ArrowRightIcon, ArrowUpIcon, MinusIcon } from 'lucide-react'
import { cn } from '@/lib/utils'
import type { Priority } from '@/data/mock'

const config: Record<Priority, { icon: React.ElementType; className: string }> = {
  urgent: { icon: AlertCircleIcon, className: 'text-red-500' },
  high: { icon: ArrowUpIcon, className: 'text-orange-500' },
  medium: { icon: ArrowRightIcon, className: 'text-yellow-500' },
  low: { icon: ArrowDownIcon, className: 'text-blue-400' },
  none: { icon: MinusIcon, className: 'text-muted-foreground' },
}

interface PriorityIconProps {
  priority: Priority
  className?: string
}

export function PriorityIcon({ priority, className }: PriorityIconProps) {
  const { icon: Icon, className: colorClass } = config[priority]
  return <Icon className={cn('size-3.5 shrink-0', colorClass, className)} />
}
