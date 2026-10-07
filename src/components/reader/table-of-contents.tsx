import { X } from "lucide-react";

interface tableOfContentsProps {
  chapters: string[]
  activeChapter: string
  open: boolean
  onClose: () => void
  onSelect: (chapter: string) => void
}

export default function TableOfContents({
  chapters,
  activeChapter,
  open,
  onClose,
  onSelect,
}: tableOfContentsProps) {
  return (
    <div
      className={`fixed top-15 bottom-0 flex w-[calc(75%-0.5rem)] flex-col overflow-y-auto bg-bg p-5 transition-all delay-75 md:sticky md:w-[calc(25%-1rem)] md:p-0 ${open ? "left-0" : "-left-full"}`}>
      <div className="mb-5 flex justify-between bg-secondary-bg p-2 text-center md:justify-center">
        <span className="font-serif font-semibold">Chapter Lists</span>
        <span
          onClick={onClose}
          className="cursor-pointer md:hidden">
          <X />
        </span>
      </div>
      <div className="flex flex-col gap-2 overflow-y-auto">
        {chapters.map((chap: string, i: number) => (
          <span
            key={i}
            className={`flex-none cursor-pointer border-b border-b-solid p-2 font-serif select-none transition-colors ${chap.toLowerCase() === activeChapter.toLowerCase()
              ? "border-b-border text-border"
              : "border-b-fg hover:border-b-border hover:text-border"}`}
            onClick={() => onSelect(chap)}>
            {chap.toUpperCase()}
          </span>
        ))}
      </div>
    </div>
  );
}