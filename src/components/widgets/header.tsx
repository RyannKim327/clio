import { House, Info, Search, User } from "lucide-react"
import { Link } from "react-router"

const navs = [
  {
    title: "Home",
    endpoint: "/",
    icon: <House />
  },
  {
    title: "About",
    endpoint: "/about",
    icon: <Info />
  },
  {
    title: "Search",
    action: "search",
    icon: <Search />
  },
  {
    title: "Login",
    action: "login",
    icon: <User />
  }
]

interface header {
  setAction: (action: string) => void
}

export default function Header({ setAction }: header) {
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
                className="list-none">
                <span className="hidden md:inline">{nav.title}</span>
                <span className="md:hidden">{nav.icon}</span>
              </li>
          )
        })}
      </nav>
    </div >
  )
}
