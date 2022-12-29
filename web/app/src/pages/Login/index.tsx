import {
  Button,
  FormControl,
  FormHelperText,
  Input,
  InputGroup,
  InputLeftElement,
  Link,
  Stack,
  Text,
} from "@chakra-ui/react";
import {AtSignIcon, LockIcon} from "@chakra-ui/icons";
import {Api, UserLogin} from "../../api"
import {ChangeEvent, FormEvent, useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../../api/useAuth";
import {AxiosError} from "axios";
import AuthLayout, {AlertData} from "../../layouts/auth";
import {Link as ReactLink, useSearchParams} from "react-router-dom";

const api = new Api();

export default function LoginPage() {
  const navigate = useNavigate();
  const {user} = useAuth({});

  const [searchParams] = useSearchParams();

  const [loginData, setLoginData] = useState<UserLogin>({username: "", password: ""});
  const [alert, setAlert] = useState<AlertData>()

  useEffect(() => {
    if (searchParams.get("registered") === "true")
      setAlert({
        status: 'success',
        title: 'Success',
        detail: 'Account created, please log in'
      })
  }, [searchParams])

  useEffect(() => {
    if (user) {
      navigate("/authenticated")
    }
  }, [user, navigate])

  function onEmailChange(e: ChangeEvent<HTMLInputElement>) {
    setLoginData({
      ...loginData,
      username: e.target.value
    })
  }

  function onPasswordChange(e: ChangeEvent<HTMLInputElement>) {
    setLoginData({
      ...loginData,
      password: e.target.value
    })
  }

  function doLogin(event: FormEvent) {
    setAlert(undefined)
    event.preventDefault()
    api.login(loginData).then(() => {
      navigate("/authenticated");
    }).catch((e: AxiosError) => {
      const err = api.parseError(e)
      console.error(err)
      setAlert({
        status: 'error',
        title: "Login failed",
        detail: err.detail
      })
    })
  }

  return (
    <AuthLayout title={"Login"} alert={alert}>
      <form onSubmit={doLogin}>
        <Stack spacing={4}>
          <FormControl>
            <InputGroup>
              <InputLeftElement
                pointerEvents="none"
                children={<AtSignIcon/>}
              />
              <Input type="email" placeholder="email address"
                     onChange={onEmailChange} required/>
            </InputGroup>
          </FormControl>
          <FormControl>
            <InputGroup>
              <InputLeftElement
                pointerEvents="none"
                children={<LockIcon/>}
              />
              <Input
                type={"password"}
                placeholder="Password"
                onChange={onPasswordChange}
                required
              />
            </InputGroup>
            <FormHelperText textAlign="right">
              <Link>forgot password?</Link>
            </FormHelperText>
          </FormControl>
          <Button
            borderRadius={0}
            type="submit"
            variant="solid"
            width="full"
          >
            Login
          </Button>
          <Text textAlign={"center"}>
            New here? <Link as={ReactLink} to={"/signup"}>Create an account</Link>
          </Text>
        </Stack>
      </form>
    </AuthLayout>
  );
}