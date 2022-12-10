import {
  Alert, AlertDescription,
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
import {AuthApi, UserLogin} from "../../api/generated";
import {ChangeEvent, useState} from "react";

export default function LoginPage() {
  const authApi = new AuthApi();

  const [loginData, setLoginData] = useState<UserLogin>({});
  const [error, setError] = useState<string>()

  function onEmailChange(e: ChangeEvent<HTMLInputElement>) {
    setLoginData({
      ...loginData,
      email: e.target.value
    })
  }

  function onPasswordChange(e: ChangeEvent<HTMLInputElement>) {
    setLoginData({
      ...loginData,
      password: e.target.value
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
          <form onSubmit={(event) => {
            setError(undefined)
            console.log(loginData)
            authApi.authLoginPost(loginData).then((result) => {
              // Do something
            }).catch((result) => {
              console.log(result)
              setError(result.response.data.message)
            })
            event.preventDefault()
          }}>
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