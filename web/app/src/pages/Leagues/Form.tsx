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
  Stack,
} from "@chakra-ui/react";
import React, {ReactNode, useState} from "react";
import {Api, ErrorResponse, LeagueCreate, useAuth} from "../../api";
import {slugify} from "../../helpers";
import {useNavigate} from "react-router-dom";
import {AxiosError} from "axios";
import Alert, {AlertData} from "../../components/Alert";
import {LocationOption, LocationSelect} from "../../components/LocationSelect";

export type LeagueFormProps = {
  mode: 'create' | 'edit'
  onCancel?: () => void
}

const pinmanApi = new Api();

export default function LeagueForm(props: LeagueFormProps) {
  const navigate = useNavigate();

  useAuth({requireAuth: true})

  const [alert, setAlert] = useState<AlertData>()
  const [locationValue, setLocationValue] = useState<LocationOption | undefined>()
  const [slugModified, setSlugModified] = useState<boolean>(false)
  const [slugError, setSlugError] = useState<ReactNode | undefined>(undefined)
  const [leagueData, setLeagueData] = useState<LeagueCreate>({
    name: "",
    location_id: "",
    slug: "",
  });

  async function submitForm() {
    let locationId: string | undefined
    if (locationValue === undefined) {
      setAlert({
        status: "error",
        title: "No location selected",
        detail: "Please select a location for your league"
      })
      return
    }
    if (locationValue.type === "pinmap") {
      locationId = await pinmanApi.locationsApi().locationsPost({
        pinball_map_id: parseInt(locationValue.pinballMapId),
      }).then((response) => response.data.location.id).catch((e: AxiosError) => {
        const err = pinmanApi.parseError(e)
        setAlert({
          status: "error",
          title: err.title,
          detail: err.detail
        })
        return undefined
      })
      if (locationId === undefined) {
        return
      }
    } else {
      locationId = locationValue.value
    }

    pinmanApi.leaguesApi().leaguesPost({
      ...leagueData,
      location_id: locationId,
    }).then(() => {
      // navigate(`/leagues/${response.data.league?.slug}`)
      navigate(`/leagues`)
    }).catch((e: AxiosError) => {
      const err = pinmanApi.parseError(e)
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
            Slug must be url-friendly. Try <code>{validSlug}</code> instead, or clear the field to use the
            generated
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

  function onLocationChange(newValue: LocationOption) {
    setLocationValue(newValue)
  }

  function onSelectError(e: ErrorResponse) {
    setAlert({
      status: "error",
      title: e.title,
      detail: e.detail
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
          <LocationSelect onChange={onLocationChange} onError={onSelectError}/>
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
