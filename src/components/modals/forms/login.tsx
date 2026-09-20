import { useState, type ChangeEvent, type SubmitEvent } from "react";

export default function LoginForm() {
  const [username, setUsername] = useState("")
  const [password, setPassword] = useState("")

  function submitForm(event: SubmitEvent<HTMLFormElement>) {

  }

  return (
    <div className="flex flex-col w-md">
      <h1 className="text-xl border-b-border border-b border-b-solid w-full text-center mb-5">Login</h1>
      <form onSubmit={submitForm}
        className="flex flex-col gap-2" >
        <div className="flex flex-col">
          <label
            htmlFor="username"
            className="text-xs">Username/Email</label>
          <input
            value={username}
            onChange={(e: ChangeEvent<HTMLInputElement, HTMLInputElement>) => {
              setUsername(e.target.value)
            }}
            className="outline-none border border-solid border-fg rounded hover:border-border px-2 py-1 text-xs"
            id="username"
            type="text" />
        </div>
        <div className="flex flex-col">
          <label
            htmlFor="password"
            className="text-xs">Password</label>
          <input
            value={password}
            onChange={(e: ChangeEvent<HTMLInputElement, HTMLInputElement>) => {
              setPassword(e.target.value)
            }}
            className="outline-none border border-solid border-fg rounded hover:border-border px-2 py-1 text-xs"
            id="password"
            type="password" />
        </div>
        <input
          className="bg-border/50 hover:bg-border/25 rounded-sm w-full text-center p-2 text-xs"
          type="submit"
          value="Login" />
      </form >
    </div >
  )
}
