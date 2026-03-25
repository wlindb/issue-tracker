import { useEffect, useRef, useState } from 'react'
import { cn } from '@/lib/utils'

interface EditableTextProps {
  value: string
  onSave: (newValue: string) => void
  placeholder?: string
  multiline?: boolean
  className?: string
  inputClassName?: string
}

export function EditableText({
  value,
  onSave,
  placeholder = 'Click to edit...',
  multiline = false,
  className,
  inputClassName,
}: EditableTextProps) {
  const [isEditing, setIsEditing] = useState(false)
  const [draft, setDraft] = useState(value)
  const inputRef = useRef<HTMLInputElement & HTMLTextAreaElement>(null)

  useEffect(() => {
    if (isEditing) {
      inputRef.current?.focus()
      inputRef.current?.select()
    }
  }, [isEditing])

  function startEditing() {
    setDraft(value)
    setIsEditing(true)
  }

  function save() {
    const trimmed = draft.trim()
    if (trimmed && trimmed !== value) {
      onSave(trimmed)
    }
    setIsEditing(false)
  }

  function discard() {
    setDraft(value)
    setIsEditing(false)
  }

  const sharedInputClass = cn(
    'w-full rounded-md border border-input bg-transparent px-2.5 py-1.5 text-sm ring-offset-background',
    'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
    inputClassName,
  )

  if (isEditing) {
    if (multiline) {
      return (
        <textarea
          ref={inputRef as React.RefObject<HTMLTextAreaElement>}
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          onBlur={save}
          onKeyDown={(e) => {
            if (e.key === 'Escape') { e.preventDefault(); discard() }
            if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); save() }
          }}
          placeholder={placeholder}
          rows={4}
          className={cn(sharedInputClass, 'resize-none leading-relaxed', className)}
        />
      )
    }

    return (
      <input
        ref={inputRef as React.RefObject<HTMLInputElement>}
        type="text"
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        onBlur={save}
        onKeyDown={(e) => {
          if (e.key === 'Escape') { e.preventDefault(); discard() }
          if (e.key === 'Enter') { e.preventDefault(); save() }
        }}
        placeholder={placeholder}
        className={cn(sharedInputClass, className)}
      />
    )
  }

  return (
    <button
      type="button"
      onClick={startEditing}
      className={cn(
        'w-full rounded-md px-2.5 py-1.5 text-left text-sm transition-colors hover:bg-muted/60',
        !value && 'text-muted-foreground italic',
        className,
      )}
    >
      {value || placeholder}
    </button>
  )
}
