import {PropsWithChildren} from "react";
import {Avatar, Box, Flex, Heading, Stack} from "@chakra-ui/react";
import ColorToggle from "../../components/ColorToggle";
import Alert, {AlertData} from "../../components/Alert";

type AuthLayoutProps = PropsWithChildren & {
  title?: string;
  alert?: AlertData
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

        {props.alert &&
          <Alert alert={props.alert}/>
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
