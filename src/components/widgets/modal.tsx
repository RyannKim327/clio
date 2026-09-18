import { X } from "lucide-react"
import type { ReactNode } from "react"

interface modal {
  show: boolean
  children: ReactNode
  closeModal: () => void
}

export default function Modal({ closeModal, show, children }: modal) {
  return (
    <div className={`${show ? "flex z-10 items-center justify-center" : "hidden"} fixed top-0 right-0 left-0 bottom-0 bg-bg/75`}>
      <span
        onClick={closeModal}
        className="cursor-pointer fixed top-5 right-5">
        <X />
      </span>
      <div className="flex max-w-2/3 max-h-2/3 bg-secondary-bg p-3 rounded">
        {children}
      </div>
    </div>
  )
}
