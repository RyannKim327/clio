import { stories } from "@/data/stories";
import { markdownToHtml } from "@/lib/markdown";
import { splitSegmentsByStartEnd } from "@/lib/split-segments";
import { ArrowLeft } from "lucide-react";
import { Link, useSearchParams } from "react-router";

export default function Read() {
  const [searchParams] = useSearchParams()
  const id = searchParams.get("id")
  const chapter = parseInt(searchParams.get("c") ?? "0")
  const story = stories[parseInt(id)]

  const chapters = splitSegmentsByStartEnd(story.chapters[chapter])

  console.log(chapters)

  return (
    <div className="flex flex-col items-center w-full h-full gap-2 overflow-x-hidden">
      <div className="flex items-center w-full bg-secondary-bg border-b-2 border-b-solid border-border sticky z-10 top-0 p-2 px-3 gap-3">
        <Link to="/">
          <ArrowLeft />
        </Link>

        <div className="flex flex-col">
          <span>{story.title}</span>
          <span className="text-xs">{story.author}</span>
        </div>
      </div>

      <div className="flex flex-col gap-24 max-w-2xl overflow-x-hidden scrollbar-thin scrollbar-thumb-secondary-bg snap-mandatory snap-y">
        <div className="sticky top-0 bg-secondary-bg p-2 text-center">
          Chapter {chapter + 1}
        </div>
        {
          chapters.map((line) => {
            return (
              <div className="min-h-full snap-start px-5 pt-15 overflow-y-auto scrollbar-none">
                {line.map((c) => {
                  return (
                    <p className="font-mono" dangerouslySetInnerHTML={{ __html: markdownToHtml(c.content) }} />
                  )
                })}
              </div>
            )
          })
        }
      </div>
    </div >
  )
}
