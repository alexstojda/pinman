import {Button} from "@chakra-ui/react";
import {Api, useAuth} from "../../api";
import {useNavigate} from "react-router-dom";

export default function AuthTest() {
  const {user} = useAuth({requireAuth: true});
  const navigate = useNavigate();

  const api = new Api()

  return (
    <>
      <h1>User is authenticated</h1>
      <pre>{JSON.stringify(user, null, 2)}</pre>
      <Button variant='outline' marginLeft={"2em"} onClick={() => {
        api.clearJwtToken()
        navigate("/login")
      }}>
        Logout
      </Button>
    </>
  )
}