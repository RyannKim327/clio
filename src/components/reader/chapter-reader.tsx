import AudioController from "@/components/reader/audio-controller";
import Paragraph from "@/components/reader/paragraph";
import type { c_chapter } from "@/types/book";
import { useCallback, useState } from "react";

interface readerState {
  chapter: string
  index: number | null
  musicUrl: string | null
  end: boolean
}

interface chapterReaderProps {
  chapter: string
  paragraphs: c_chapter[]
  canPlay: boolean
  volume: number
  muted: boolean
  onPlayingChange: (playing: boolean) => void
}

const idle = (chapter: string, end: boolean): readerState => ({
  chapter,
  index: null,
  musicUrl: null,
  end,
});

export default function ChapterReader({
  chapter,
  paragraphs,
  canPlay,
  volume,
  muted,
  onPlayingChange,
}: chapterReaderProps) {
  const [active, setActive] = useState<readerState>(() => idle(chapter, false));

  // A new chapter must stop any audio that was playing. React re-renders
  // synchronously after this adjustment, before the new paragraphs' observers fire.
  if (active.chapter !== chapter) {
    setActive(idle(chapter, true));
  }

  const handleViewportChange = useCallback((index: number, paragraph: c_chapter) => {
    setActive((prev) => {
      // `end` takes precedence over `music` per the data contract.
      if (paragraph.end) return { ...prev, index, musicUrl: null, end: true };
      if (paragraph.music) return { ...prev, index, musicUrl: paragraph.music, end: false };
      // Plain paragraphs never change the audio; the current track continues.
      return { ...prev, index };
    });
  }, []);

  return (
    <>
      <AudioController
        activeMusicUrl={active.musicUrl}
        shouldStop={active.end}
        volume={volume}
        muted={muted}
        canPlay={canPlay}
        onPlayingChange={onPlayingChange}
      />
      <div className="flex flex-col gap-6">
        {paragraphs.map((paragraph, i) => (
          <Paragraph
            key={`${chapter}:${i}`}
            index={i}
            content={paragraph.content}
            music={paragraph.music}
            end={paragraph.end}
            isActive={active.index === i}
            onViewportChange={handleViewportChange}
          />
        ))}
      </div>
    </>
  );
}