import {useAuth} from "../../api/useAuth";

export default function AuthTest() {
  const [user] = useAuth({requireAuth: true});

  return (
    <>
      <h1>User is authenticated</h1>
      <pre>{JSON.stringify(user, null, 2)}</pre>
    </>
  )
}