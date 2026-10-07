import ChapterReader from "@/components/reader/chapter-reader";
import NowPlaying from "@/components/reader/now-playing";
import StartOverlay from "@/components/reader/start-overlay";
import TableOfContents from "@/components/reader/table-of-contents";
import { stories } from "@/data/stories";
import type { c_chapter } from "@/types/book";
import { ArrowLeft, Menu } from "lucide-react";
import { useCallback, useLayoutEffect, useRef, useState } from "react";
import { Link, useLocation, useNavigate, useParams, useSearchParams } from "react-router";

export default function Read() {
  const [searchParams] = useSearchParams()
  const { id } = useParams()
  const location = useLocation()
  const navigate = useNavigate()
  const chapter = searchParams.get("c") ?? "chapter 1"
  const story = stories[parseInt(id ?? "0")]
  const paragraphs: c_chapter[] = story.chapters[chapter.toLowerCase()] ?? []
  const scrollRef = useRef<HTMLDivElement | null>(null)

  const [showChapters, setShowChapters] = useState(false)
  const [hasInteracted, setHasInteracted] = useState(false)
  const [isPlaying, setIsPlaying] = useState(false)
  const [muted, setMuted] = useState(false)
  const [volume] = useState(0.85)

  const handlePlayingChange = useCallback((playing: boolean) => {
    setIsPlaying(playing)
  }, [])

  // Jump to the top whenever the selected chapter changes so the reader
  // starts at the chapter heading. Runs before paint so the new paragraphs'
  // IntersectionObserver entries are computed against the reset position.
  useLayoutEffect(() => {
    scrollRef.current?.scrollTo({ top: 0 })
  }, [chapter])

  function handleSelectChapter(chap: string) {
    setShowChapters(false)
    const params = new URLSearchParams(location.search)
    params.set("c", chap)
    navigate(`${location.pathname}?${params.toString()}`)
  }

  return (
    <div className="flex flex-col items-center w-full h-full gap-2 overflow-x-hidden">
      {!hasInteracted && (
        <StartOverlay
          title={story.title}
          description={story.description}
          onStart={() => setHasInteracted(true)}
        />
      )}

      <div className="flex items-center w-full bg-secondary-bg border-b-2 border-b-solid border-border sticky z-10 top-0 p-2 px-3 gap-3">
        <Link to="/">
          <ArrowLeft />
        </Link>

        <div className="flex flex-col">
          <span>{story.title}</span>
          <span className="text-xs">{story.author}</span>
        </div>
      </div>

      <div className="flex w-[calc(90%-1rem)] justify-between h-full overflow-hidden">
        <TableOfContents
          chapters={Object.keys(story.chapters)}
          activeChapter={chapter}
          open={showChapters}
          onClose={() => setShowChapters(false)}
          onSelect={handleSelectChapter}
        />
        <div className="flex flex-col w-full md:w-[calc(70%-1rem)] gap-2">
          <div className="flex gap-2 md:justify-center bg-secondary-bg p-2 text-center font-serif font-semibold">
            <span
              onClick={() => {
                setShowChapters(true)
              }}
              className="md:hidden cursor-pointer">
              <Menu />
            </span>
            {chapter.toUpperCase()}
          </div>
          <div
            ref={scrollRef}
            className="flex flex-col overflow-x-hidden scrollbar-thin scrollbar-thumb-secondary-bg p-5"
          >
            <ChapterReader
              chapter={chapter}
              paragraphs={paragraphs}
              canPlay={hasInteracted}
              volume={volume}
              muted={muted}
              onPlayingChange={handlePlayingChange}
            />
          </div>
        </div>
      </div>

      <NowPlaying
        visible={isPlaying && hasInteracted}
        muted={muted}
        onToggleMute={() => setMuted((prev) => !prev)}
      />
    </div >
  )
}