import {
  Button,
  FormControl,
  Input,
  InputGroup,
  Link,
  Stack,
  Text,
} from "@chakra-ui/react";
import {Api, UserRegister} from "../../api"
import {ChangeEvent, FormEvent, useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../../api/useAuth";
import {AxiosError} from "axios";
import AuthLayout, {AlertData} from "../../layouts/auth";
import {Link as ReactLink} from "react-router-dom";

const api = new Api();

export default function SignUpPage() {
  const navigate = useNavigate();

  const [user] = useAuth({});

  const [registrationData, setRegistrationData] = useState<UserRegister>({
    name: "",
    email: "",
    password: "",
    passwordConfirm: "",
  });
  const [alert, setAlert] = useState<AlertData>()

  useEffect(() => {
    if (user) {
      navigate("/authenticated")
    }
  }, [user, navigate])

  function onNameChange(e: ChangeEvent<HTMLInputElement>) {
    setRegistrationData({
      ...registrationData,
      name: e.target.value
    })
  }

  function onEmailChange(e: ChangeEvent<HTMLInputElement>) {
    setRegistrationData({
      ...registrationData,
      email: e.target.value
    })
  }

  function onPasswordChange(e: ChangeEvent<HTMLInputElement>) {
    setRegistrationData({
      ...registrationData,
      password: e.target.value
    })
  }

  function onPasswordConfirmChange(e: ChangeEvent<HTMLInputElement>) {
    setRegistrationData({
      ...registrationData,
      passwordConfirm: e.target.value
    })
  }

  function doRegisterUser(event: FormEvent) {
    setAlert(undefined)
    event.preventDefault()

    if (registrationData.password !== registrationData.passwordConfirm) {
      setAlert({
        status: "error",
        title: "Error",
        detail: "Password confirmation does not match"
      })
      return
    }

    api.userApi().usersRegisterPost(registrationData).then(() => {
      navigate("/login?registered=true");
    }).catch((e: AxiosError) => {
      const err = api.parseError(e)
      console.error(err)
      setAlert({
        status: "error",
        title: err.title,
        detail: err.detail
      })
    })
  }

  return (
    <AuthLayout title={"Register"} alert={alert}>
      <form onSubmit={doRegisterUser}>
        <Stack spacing={4}>
          <FormControl>
            <InputGroup>
              <Input type="text" placeholder="name"
                     onChange={onNameChange} required/>
            </InputGroup>
          </FormControl>
          <FormControl>
            <InputGroup>
              <Input type="email" placeholder="email address"
                     onChange={onEmailChange} required/>
            </InputGroup>
          </FormControl>
          <FormControl>
            <InputGroup>
              <Input
                type={"password"}
                placeholder="password"
                onChange={onPasswordChange}
                required
              />
            </InputGroup>
          </FormControl>
          <FormControl>
            <InputGroup>
              <Input
                type={"password"}
                placeholder="confirm password"
                onChange={onPasswordConfirmChange}
                required
              />
            </InputGroup>
          </FormControl>
          <Button
            borderRadius={0}
            type="submit"
            variant="solid"
            width="full"
          >
            Create account
          </Button>
          <Text textAlign={"center"}>
            Already have an account? <Link as={ReactLink} to={"/login"}>Sign in</Link>
          </Text>
        </Stack>
      </form>
    </AuthLayout>
  );
}