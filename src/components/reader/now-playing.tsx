import { Music, Volume2, VolumeX } from "lucide-react";

interface nowPlayingProps {
  visible: boolean
  muted: boolean
  onToggleMute: () => void
}

export default function NowPlaying({ visible, muted, onToggleMute }: nowPlayingProps) {
  if (!visible) return null;

  return (
    <div
      role="status"
      aria-live="polite"
      className="fixed right-4 bottom-4 z-40 flex items-center gap-2 rounded-full border border-border bg-secondary-bg px-3 py-2 select-none">
      <Music size={14} className="animate-pulse text-border" aria-hidden />
      <span className="text-xs font-serif">Now Playing</span>
      <button
        onClick={onToggleMute}
        aria-label={muted ? "Unmute" : "Mute"}
        className="cursor-pointer rounded-full p-1 transition-colors hover:bg-border/20">
        {muted ? <VolumeX size={14} aria-hidden /> : <Volume2 size={14} aria-hidden />}
      </button>
    </div>
  );
}