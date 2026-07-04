import { useRef, useState } from 'react'
import { PlusIcon } from 'lucide-react'
import { createLabel, type Label } from '@/api/generated/issueTrackerAPI'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Checkbox } from '@/components/ui/checkbox'
import { LabelDot } from '@/components/LabelDot'
import { useLabelSearch } from '@/hooks/useLabelSearch'

interface LabelPickerProps {
  workspaceId: string
  attachedLabels: Label[]
  onLabelsChange: (labels: Label[]) => void
}

export function LabelPicker({ workspaceId, attachedLabels, onLabelsChange }: LabelPickerProps) {
  const [query, setQuery] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)
  const { results, isPending } = useLabelSearch(workspaceId, query)

  function toggleLabel(label: Label) {
    const isAttached = attachedLabels.some((l) => l.id === label.id)
    if (isAttached) {
      onLabelsChange(attachedLabels.filter((l) => l.id !== label.id))
    } else {
      onLabelsChange([...attachedLabels, label])
    }
  }

  async function handleCreateLabel() {
    const name = query.trim()
    if (!name) return
    const created = await createLabel(workspaceId, { name })
    onLabelsChange([...attachedLabels, created])
    setQuery('')
  }

  const trimmedQuery = query.trim()
  const hasExactMatch = results.some((l) => l.name.toLowerCase() === trimmedQuery.toLowerCase())
  const showCreateRow = trimmedQuery !== '' && !isPending && !hasExactMatch

  return (
    <Popover>
      <PopoverTrigger
        render={
          <Button variant="ghost" size="icon-xs" aria-label="Add label">
            <PlusIcon />
          </Button>
        }
      />
      <PopoverContent initialFocus={inputRef} className="p-2">
        <Input
          ref={inputRef}
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Add labels..."
          className="mb-1.5"
        />
        <div className="max-h-64 overflow-y-auto">
          {results.map((label) => {
            const isAttached = attachedLabels.some((l) => l.id === label.id)
            return (
              <label
                key={label.id}
                className="flex cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-sm hover:bg-muted"
              >
                <Checkbox checked={isAttached} onCheckedChange={() => toggleLabel(label)} />
                <LabelDot labelId={label.id} />
                <span>{label.name}</span>
              </label>
            )
          })}
          {showCreateRow && (
            <button
              type="button"
              onClick={handleCreateLabel}
              className="flex w-full cursor-pointer items-center gap-2 rounded-md px-2 py-1.5 text-left text-sm hover:bg-muted"
            >
              <PlusIcon className="size-3.5 shrink-0 text-muted-foreground" />
              <span>
                Create new workspace label: <span className="font-medium">&quot;{trimmedQuery}&quot;</span>
              </span>
            </button>
          )}
        </div>
      </PopoverContent>
    </Popover>
  )
}
