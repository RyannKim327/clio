import { stories } from "@/data/stories";
import { markdownToHtml } from "@/lib/markdown";
import { splitSegmentsByStartEnd } from "@/lib/split-segments";
import { ArrowLeft, Menu, X } from "lucide-react";
import { useState } from "react";
import { Link, useLocation, useNavigate, useParams, useSearchParams } from "react-router";

export default function Read() {
  const [searchParams] = useSearchParams()
  const { id } = useParams()
  const location = useLocation()
  const navigate = useNavigate()
  const chapter = searchParams.get("c") ?? "chapter 1"
  const story = stories[parseInt(id)]
  const [showChapters, setShowChapters] = useState(false)
  const chapters = splitSegmentsByStartEnd(story.chapters[chapter.toLowerCase()])

  function toTop() {
    const _ = document.getElementById("headToTop");
    if (_) {
      _.scrollIntoView({
        behavior: "smooth",
      });
    }
  }

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

      <div className="flex w-[calc(90%-1rem)] justify-between h-full overflow-hidden">
        <div className={`flex flex-col w-full md:w-[calc(25%-1rem)] overflow-y-auto fixed ${showChapters ? "left-0" : "-left-full"} p-5 md:p-0 top-15 bottom-0 bg-bg md:sticky transition-all delay-75`}>
          <div className="flex justify-between md:justify-center bg-secondary-bg p-2 text-center mb-5">
            <span>Chapter Lists</span>
            <span
              onClick={() => {
                setShowChapters(false)
              }}
              className={`md:hidden`}>
              <X />
            </span>
          </div>
          <div className="flex flex-col gap-2 overflow-y-auto">
            {
              Object.keys(story.chapters).map((chap, i: number) => {
                return (
                  <span
                    key={i}
                    className={`border-b border-b-solid border-fg hover:border-border hover:text-border p-2 cursor-pointer select-none ${chap.toLowerCase() === chapter.toLowerCase() ? "border-b-border text-border" : ""}`}
                    onClick={() => {
                      toTop()
                      setShowChapters(false)
                      const params = new URLSearchParams(location.search)
                      params.set("c", chap)
                      navigate(`${location.pathname}?${params.toString()}`)
                    }}
                  > {chap.toUpperCase()}</span>
                )
              })
            }
          </div>
        </div>
        <div className="flex flex-col w-full md:w-[calc(70%-1rem)] gap-2">
          <div className="flex gap-2 md:justify-center bg-secondary-bg p-2 text-center">
            <span
              onClick={() => {
                setShowChapters(true)
              }}
              className="md:hidden">
              <Menu />
            </span>
            {chapter.toUpperCase()}
          </div>
          <div
            className="flex flex-col overflow-x-hidden scrollbar-thin scrollbar-thumb-secondary-bg snap-mandatory snap-y p-5"
          >
            {
              chapters.map((line, i: number) => {
                return (
                  <div
                    key={i}
                    className="min-h-full snap-start overflow-y-auto scrollbar-none">

                    {line.map((c, j: number) => {
                      return (
                        <p
                          id={`${i}:${j}`}
                          key={`${i}:${j}`}
                          className="font-semibold font-serif"
                          dangerouslySetInnerHTML={{ __html: markdownToHtml(c.content) }} />
                      )
                    })}
                  </div>
                )
              })
            }
          </div>
        </div>
      </div>
    </div >
  )
}
