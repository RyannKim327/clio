import { X } from "lucide-react"
import type { MouseEvent, ReactNode } from "react"

interface modal {
  show: boolean
  children: ReactNode
  closeModal: () => void
  className?: string
}

export default function Modal({ className, closeModal, show, children }: modal) {
  return (
    <div
      onClick={(e: MouseEvent<HTMLDivElement, globalThis.MouseEvent>) => {
        if (e.target === e.currentTarget) {
          closeModal()
        }
      }}
      className={`${show ? "flex z-10 items-center justify-center" : "hidden"} fixed top-0 right-0 left-0 bottom-0 bg-bg/75`}>
      <span
        onClick={() => {
          closeModal()
        }} className="cursor-pointer fixed top-5 right-5 pointer-events-auto">
        <X />
      </span>
      <div
        className={`${className ?? "max-w-[calc(75%-0.5rem)] md:max-w-[calc(50%-1rem)]"} select-none max-h-[calc(50%-0.5rem)] flex bg-secondary-bg p-3 rounded pointer-events-auto overflow-y-auto scrollbar-thin scrollbar-thumb-border scrollbar-track-bg`}>
        {children}
      </div>
    </div>
  )
}
