import { useParagraphObserver } from "@/hooks/use-paragraph-observer";
import { markdownToHtml } from "@/lib/markdown";
import type { c_chapter } from "@/types/book";
import { Music, Square } from "lucide-react";
import { useEffect } from "react";

interface paragraphProps {
  index: number
  content: string
  music?: string
  end?: boolean
  isActive: boolean
  onViewportChange: (index: number, paragraph: c_chapter) => void
}

export default function Paragraph({
  index,
  content,
  music,
  end,
  isActive,
  onViewportChange,
}: paragraphProps) {
  const [ref, isIntersecting] = useParagraphObserver<HTMLDivElement>(0.5);

  useEffect(() => {
    if (isIntersecting) {
      onViewportChange(index, { content, music, end });
    }
  }, [isIntersecting, index, content, music, end, onViewportChange]);

  return (
    <div
      ref={ref}
      className={`relative border-l-[3px] pl-4 transition-colors ${isActive ? "border-border bg-secondary-bg" : "border-transparent"}`}>
      {music && (
        <span
          className="flex items-center gap-1 text-xs text-border/80"
          title="Audio starts here">
          <Music size={12} aria-hidden />
          <span className="sr-only">Audio starts here</span>
        </span>
      )}
      <div
        className="font-serif"
        dangerouslySetInnerHTML={{ __html: markdownToHtml(content) }}
      />
      {end && (
        <div
          className="mt-6 mb-4 flex items-center gap-1 text-xs text-border/80 select-none"
          aria-label="Audio ends here">
          <Square size={10} aria-hidden />
          <span>Audio stops here</span>
        </div>
      )}
    </div>
  );
}