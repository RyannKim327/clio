import { BookOpen } from "lucide-react";

interface startOverlayProps {
  title: string
  description?: string
  onStart: () => void
}

export default function StartOverlay({ title, description, onStart }: startOverlayProps) {
  return (
    <div
      role="dialog"
      aria-modal="true"
      aria-label="Start reading"
      className="fixed inset-0 z-50 flex flex-col items-center justify-center gap-4 bg-bg/95 p-6 text-center backdrop-blur-sm">
      <span className="flex items-center gap-2 text-border">
        <BookOpen size={20} aria-hidden />
      </span>
      <h1 className="max-w-md font-serif text-2xl leading-snug md:text-3xl">{title}</h1>
      {description && (
        <p className="max-w-md font-serif text-sm leading-relaxed opacity-75">{description}</p>
      )}
      <button
        onClick={onStart}
        className="mt-4 cursor-pointer border border-border px-8 py-2 font-serif transition-colors hover:bg-border hover:text-bg">
        Start Reading
      </button>
      <span className="text-xs opacity-60">Enable sound to let the story play its music</span>
    </div>
  );
}