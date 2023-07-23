import {Alert as ChakraAlert, AlertDescription, AlertIcon, AlertTitle} from "@chakra-ui/react";

export type AlertData = {
  status: 'error' | 'success';
  title: string;
  detail: string;
}

type AlertProps = {
  alert: AlertData
}

export default function Alert(props: AlertProps) {
  return (
    <ChakraAlert status={props.alert.status}>
      <AlertIcon/>
      <AlertTitle>{props.alert.title}</AlertTitle>
      <AlertDescription>{props.alert.detail}</AlertDescription>
    </ChakraAlert>
  )
}
