import ColorToggle from "../../components/ColorToggle";
import {Text, HStack, Spacer, Link} from "@chakra-ui/react";

export default function Footer() {
  return (
    <HStack my={4} >
      <Spacer/>
      <Text fontSize="sm" color="gray.500">
        © {new Date().getFullYear()} - PinMan | <Link href={"https://github.com/alexstojda/pinman"}>GitHub</Link>
      </Text>
      <ColorToggle/>
      <Spacer/>
    </HStack>
  );
}
