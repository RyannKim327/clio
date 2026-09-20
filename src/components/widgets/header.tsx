import { House, Info, Search as s, User } from "lucide-react"
import { Link } from "react-router"
import About from "@/components/modals/about"
import Search from "@/components/modals/forms/search"
import Modal from "./modal"
import Login from "@/components/modals/forms/login"

const navs = [
  {
    title: "Home",
    endpoint: "/",
    icon: <House />
  },
  {
    title: "About",
    action: "about",
    icon: <Info />
  },
  {
    title: "Search",
    action: "search",
    icon: <s />
  },
  {
    title: "Login",
    action: "login",
    icon: <User />
  }
]

interface header {
  setAction: (action: string) => void
  action: string
}

export default function Header({ setAction, action }: header) {

  return (
    <div className="flex justify-between items-center w-full h-10 bg-secondary-bg border-b-2 border-b-solid border-border px-2 sticky top-0 z-10 mb-2">
      <span>Clio</span>

      <nav className="flex gap-2">
        {navs.map((nav, i: number) => {
          return (
            nav.endpoint ?
              <Link
                key={`${i}. ${nav.title}`}
                to={nav.endpoint}>
                <span className="hidden md:inline">{nav.title}</span>
                <span className="md:hidden">{nav.icon}</span>
              </Link> :
              <li
                onClick={() => {
                  setAction(nav.action ?? "")
                }}
                key={`${i}. ${nav.title}`}
                className="list-none cursor-pointer">
                <span className="hidden md:inline">{nav.title}</span>
                <span className="md:hidden">{nav.icon}</span>
              </li>
          )
        })}
      </nav>
      <Modal
        className="max-w-[calc(50%-1rem)]"
        closeModal={() => { setAction("") }}
        show={action === "about"}>
        <About />
      </Modal>
      <Modal
        closeModal={() => { setAction("") }}
        show={action === "search"}>
        <Search />
      </Modal>
      <Modal
        closeModal={() => { setAction("") }}
        show={action === "login"}>
        <Login />
      </Modal>
    </div>
  )
}
