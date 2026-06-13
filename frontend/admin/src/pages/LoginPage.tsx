export function LoginPage() {
  return (
    <div>
      <h1>Admin Login</h1>
      <form>
        <label htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
        />
        <label htmlFor="password">Password</label>
        <input
          id="password"
          type="password"
        />
        <button type="submit">Login</button>
      </form>
    </div>
  )
}
