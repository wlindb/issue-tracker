import { useState } from 'react'
import { SendHorizontalIcon } from 'lucide-react'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import type { Comment } from '@/api/generated/issueTrackerAPI'

interface CommentSectionProps {
  comments: Comment[]
  onAddComment: (body: string) => void
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

export function CommentSection({ comments, onAddComment }: CommentSectionProps) {
  const [body, setBody] = useState('')

  function handleSubmit() {
    const trimmed = body.trim()
    if (!trimmed) return
    onAddComment(trimmed)
    setBody('')
  }

  return (
    <div className="flex flex-col gap-4">
      <h3 className="text-sm font-semibold">Comments</h3>

      {comments.length > 0 && (
        <div className="flex flex-col gap-4">
          {comments.map((comment, i) => (
            <div key={comment.id}>
              {i > 0 && <Separator className="mb-4" />}
              <CommentItem comment={comment} />
            </div>
          ))}
        </div>
      )}

      <div className="flex gap-2 rounded-lg border border-border p-3">
        <textarea
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="Leave a comment..."
          rows={2}
          className="min-h-0 flex-1 resize-none bg-transparent text-sm placeholder:text-muted-foreground focus:outline-none"
          onKeyDown={(e) => {
            if (e.key === 'Enter' && e.ctrlKey) { e.preventDefault(); handleSubmit() }
          }}
        />
        <Button
          size="icon"
          variant="ghost"
          onClick={handleSubmit}
          disabled={!body.trim()}
          className="self-end shrink-0"
        >
          <SendHorizontalIcon className="size-4" />
        </Button>
      </div>
    </div>
  )
}

function CommentItem({ comment }: { comment: Comment }) {
  const shortId = comment.authorId.slice(0, 8)
  const initials = comment.authorId.slice(0, 2).toUpperCase()
  return (
    <div className="flex gap-3">
      <Avatar size="sm" className="mt-0.5 shrink-0">
        <AvatarFallback>{initials}</AvatarFallback>
      </Avatar>
      <div className="flex flex-col gap-1">
        <div className="flex items-baseline gap-2">
          <span className="text-sm font-medium font-mono">{shortId}</span>
          <span className="text-xs text-muted-foreground">{formatDate(comment.createdAt)}</span>
        </div>
        <p className="text-sm leading-relaxed text-foreground">{comment.body}</p>
      </div>
    </div>
  )
}
