import {
  Alert,
  AlertDescription,
  AlertIcon,
  AlertTitle,
  Avatar,
  Box,
  Button,
  Flex,
  FormControl,
  FormHelperText,
  Heading,
  Input,
  InputGroup,
  InputLeftElement,
  Link,
  Stack,
} from "@chakra-ui/react";
import ColorToggle from "../../components/ColorToggle";
import {AtSignIcon, LockIcon} from "@chakra-ui/icons";
import {Api, UserLogin} from "../../api"
import {ChangeEvent, FormEvent, useEffect, useState} from "react";
import {useNavigate} from "react-router-dom";
import {useAuth} from "../../api/useAuth";
import {AxiosError} from "axios";

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
    <Flex
      flexDirection="column"
      width="100wh"
      height="100vh"
      justifyContent="center"
      alignItems="center"
    >
      <Stack
        flexDir="column"
        mb="2"
        justifyContent="center"
        alignItems="center"
      >
        <Avatar/>
        <Heading>Welcome</Heading>

        <Box minW={{base: "90%", md: "468px"}}>
          <form onSubmit={doLogin}>
            <Stack
              spacing={4}
              p="1rem"
              boxShadow="md"
            >
              {error &&
                <Alert status='error'>
                  <AlertIcon/>
                  <AlertTitle>Login Failed</AlertTitle>
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              }
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
        </Box>
      </Stack>
      <Box>
        New to us?{" "}
        <Link href="#">
          Sign Up
        </Link>
        &nbsp;|&nbsp;
        <ColorToggle/>
      </Box>
    </Flex>
  );
}