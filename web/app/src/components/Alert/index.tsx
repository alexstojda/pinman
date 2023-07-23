import {Alert as ChakraAlert, AlertDescription, AlertIcon, AlertTitle} from "@chakra-ui/react";

export type AlertData = {
  status: 'error' | 'success';
  title: string;
  detail: string;
}

export default function Alert(props: AlertData) {
  return (
    <ChakraAlert status={props.status}>
      <AlertIcon/>
      <AlertTitle>{props.title}</AlertTitle>
      <AlertDescription>{props.detail}</AlertDescription>
    </ChakraAlert>
  )
}
