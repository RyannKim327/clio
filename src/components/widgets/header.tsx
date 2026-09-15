import { Link } from "react-router"

const navs = [
  {
    title: "Home",
    endpoint: ""
  },
  {
    title: "About",
    endpoint: "about"
  },
  {
    title: "Search",
    action: "search"
  },
  {
    title: "Login",
    action: "login"
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
                to={nav.endpoint}>{nav.title}</Link> :
              <li
                onClick={() => {
                  setAction(nav.action ?? "")
                }}
                key={`${i}. ${nav.title}`}
                className="list-none">{nav.title}</li>
          )
        })}
      </nav>
    </div>
  )
}
