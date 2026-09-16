import { stories } from "@/data/stories";
import { markdownToHtml } from "@/lib/markdown";
import { ArrowLeft } from "lucide-react";
import { Link, useSearchParams } from "react-router";

export default function Read() {
  const [searchParams] = useSearchParams()
  const id = searchParams.get("id")

  const story = stories[parseInt(id)]

  return (
    <div className="flex flex-col items-center w-full h-full gap-2 overflow-x-hidden">
      <div className="flex items-center w-full bg-secondary-bg border-b-2 border-b-solid border-border sticky z-10 top-0 px-3 gap-3">
        <Link to="/">
          <ArrowLeft />
        </Link>
        <div className="flex flex-col">
          <span>{story.title}</span>
          <span className="text-xs">{story.author}</span>
        </div>
      </div>
      <div className="flex flex-col gap-2 max-w-2xl overflow-x-hidden scrollbar-thin scrollbar-thumb-secondary-bg">
        {
          story.storyline.map((line) => {
            return (
              <p className="font-mono" dangerouslySetInnerHTML={{ __html: markdownToHtml(line.content) }} />
            )
          })
        }
      </div>
    </div >
  )
}
