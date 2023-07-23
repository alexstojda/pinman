import {
  Button,
  FormControl,
  FormErrorMessage,
  FormHelperText,
  FormLabel,
  HStack,
  Input,
  InputGroup,
  Spacer,
  Stack
} from "@chakra-ui/react";
import React, {ReactNode, useState} from "react";
import {Api, LeagueCreate, useAuth} from "../../api";
import {slugify} from "../../helpers";
import {useNavigate} from "react-router-dom";
import {AxiosError} from "axios";
import Alert, {AlertData} from "../../components/Alert";

export type LeagueFormProps = {
  mode: 'create' | 'edit'
  onCancel?: () => void
}

const api = new Api();

export default function LeagueForm(props: LeagueFormProps) {
  const navigate = useNavigate();

  useAuth({requireAuth: true})

  const [alert, setAlert] = useState<AlertData>()

  const [slugModified, setSlugModified] = useState<boolean>(false)
  const [slugError, setSlugError] = useState<ReactNode | undefined>(undefined)

  const [leagueData, setLeagueData] = useState<LeagueCreate>({
    name: "",
    location: "",
    slug: "",
  });

  function submitForm() {
    api.leaguesApi().leaguesPost(leagueData).then(() => {
      // navigate(`/leagues/${response.data.league?.slug}`)
      navigate(`/leagues`)
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

  function onNameChange(e: React.ChangeEvent<HTMLInputElement>) {
    setLeagueData({
      ...leagueData,
      name: e.target.value,
      ...(slugModified ? {} : {
        slug: slugify(e.target.value)
      })
    })
  }

  function onSlugChange(e: React.ChangeEvent<HTMLInputElement>) {
    setSlugError(undefined)

    if (e.target.value === "") {
      setSlugModified(false)
      setLeagueData({
        ...leagueData,
        slug: slugify(leagueData.name)
      })
    } else {
      setSlugModified(true)
      const validSlug = slugify(e.target.value)

      if (validSlug !== e.target.value) {
        setSlugError(
          <p>
            Slug must be url-friendly. Try <code>{validSlug}</code> instead, or clear the field to use the generated
            slug.
          </p>
        )
      }

      setLeagueData({
        ...leagueData,
        slug: e.target.value
      })
    }
  }

  function onLocationChange(e: React.ChangeEvent<HTMLInputElement>) {
    setLeagueData({
      ...leagueData,
      location: e.target.value
    })
  }

  return (
    <>
      {alert &&
        <Alert alert={alert}/>
      }
      <Stack spacing={4}>
        <FormControl isRequired>
          <FormLabel>League Name</FormLabel>
          <InputGroup>
            <Input type="text" placeholder="name"
                   onChange={onNameChange} required/>
          </InputGroup>
        </FormControl>
        <FormControl isRequired isInvalid={slugError !== undefined}>
          <FormLabel>Slug</FormLabel>
          <InputGroup>
            <Input type="text" placeholder="slug" value={leagueData.slug}
                   onChange={onSlugChange} required/>
          </InputGroup>
          {slugError &&
            <FormErrorMessage>{slugError}</FormErrorMessage>
          }
          <FormHelperText>
            A url-friendly identifier for your league. Use the generated slug or create your own
          </FormHelperText>
        </FormControl>
        <FormControl isRequired>
          <FormLabel>Location</FormLabel>
          <InputGroup>
            <Input type="text" placeholder="location"
                   onChange={onLocationChange} required/>
          </InputGroup>
          <FormHelperText>
            The location where the league takes place
          </FormHelperText>
        </FormControl>
      </Stack>
      <HStack my={4}>
        <Spacer/>
        <Button colorScheme='blue' mr={3} onClick={submitForm}>
          Create
        </Button>
        {props.onCancel &&
          <Button variant='ghost' onClick={props.onCancel}>
            Cancel
          </Button>
        }
      </HStack>
    </>
  )
}
