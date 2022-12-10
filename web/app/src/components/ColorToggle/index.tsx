import {Button, useColorMode} from "@chakra-ui/react";
import {MoonIcon, SunIcon} from "@chakra-ui/icons"

export default function ColorToggle() {
  const {colorMode, toggleColorMode} = useColorMode()
  return (
    <Button onClick={toggleColorMode}>
      {colorMode === 'light' ? (<SunIcon/>) : (<MoonIcon/>)}
    </Button>
  )
}