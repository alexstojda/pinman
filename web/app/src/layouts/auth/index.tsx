import {PropsWithChildren} from "react";
import {Alert, AlertDescription, AlertIcon, AlertTitle, Avatar, Box, Flex, Heading, Stack} from "@chakra-ui/react";
import ColorToggle from "../../components/ColorToggle";

type AuthLayoutProps = PropsWithChildren & {
  title?: string;
  error?: {
    title: string;
    detail: string;
  }
}

export default function AuthLayout(props: AuthLayoutProps) {
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
        width={{base: "90%", md: "60%", lg: "40%"}}
      >
        <Avatar/>
        {props.title &&
          <Heading>{props.title}</Heading>
        }

        {props.error &&
          <Alert status='error'>
            <AlertIcon/>
            <AlertTitle>{props.error.title}</AlertTitle>
            <AlertDescription>{props.error.detail}</AlertDescription>
          </Alert>
        }

        <Box
          p="1rem"
          boxShadow="md"
          width="100%"
        >
          {props.children}
        </Box>
      </Stack>
      <Box>
        <ColorToggle/>
      </Box>
    </Flex>
  )
}