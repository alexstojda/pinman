import {Button, FormControl, FormHelperText, Input, InputGroup, InputLeftElement, Link, Stack,} from "@chakra-ui/react";
import {AtSignIcon, LockIcon} from "@chakra-ui/icons";
import {Api, UserLogin} from "../../api"
import {ChangeEvent, FormEvent, useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../../api/useAuth";
import {AxiosError} from "axios";
import AuthLayout from "../../layouts/auth";

const api = new Api();

export default function LoginPage() {
  const navigate = useNavigate();

  const [user] = useAuth({});

  const [loginData, setLoginData] = useState<UserLogin>({username: "", password: ""});
  const [error, setError] = useState<string>()

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
    setError(undefined)
    event.preventDefault()
    api.login(loginData).then(() => {
      navigate("/authenticated");
    }).catch((e: AxiosError) => {
      const err = api.parseError(e)
      console.error(err)
      setError(`${err.title}: ${err.detail}`)
    })
  }

  return (
    <AuthLayout title={"Login"} error={error ? {title: "Login Failed", detail: error} : undefined}>
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
        </Stack>
      </form>
    </AuthLayout>
  );
}