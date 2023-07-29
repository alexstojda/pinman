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
import React, {ReactNode, useEffect, useState} from "react";
import {Api, LeagueCreate, pinballMap, useAuth} from "../../api";
import {slugify} from "../../helpers";
import {useNavigate} from "react-router-dom";
import {AxiosError} from "axios";
import Alert, {AlertData} from "../../components/Alert";
import {ActionMeta, Select} from "chakra-react-select";

export type LeagueFormProps = {
  mode: 'create' | 'edit'
  onCancel?: () => void
}

const pinmanApi = new Api();
const pinballMapApi = new pinballMap.Api();

export default function LeagueForm(props: LeagueFormProps) {
  const navigate = useNavigate();

  useAuth({requireAuth: true})

  const [alert, setAlert] = useState<AlertData>()

  const [slugModified, setSlugModified] = useState<boolean>(false)
  const [slugError, setSlugError] = useState<ReactNode | undefined>(undefined)

  const [leagueData, setLeagueData] = useState<LeagueCreate>({
    name: "",
    locationId: "",
    slug: "",
  });


  async function submitForm() {
    let locationId = ""
    if (locationValue!.type === "pinmap") {
      locationId = await pinmanApi.locationsApi().locationsPost({
        pinball_map_id: parseInt(locationValue!.pinballMapId),
      }).then((response) => response.data.location!.id).catch((e: AxiosError) => {
        const err = pinmanApi.parseError(e)
        setAlert({
          status: "error",
          title: err.title,
          detail: err.detail
        })
        throw new Error(err.detail)
      })
    } else {
      locationId = locationValue!.value
    }

    pinmanApi.leaguesApi().leaguesPost({
      ...leagueData,
      locationId: locationId,
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

  function onLocationChange(newValue: any, actionMeta: ActionMeta<any>) {
    setLocationValue(newValue)
  }

  type LocationOption = {
    label: string
    value: string
    type: "pinman" | "pinmap"
    pinballMapId: string
  }

  type LocationOptionGroup = {
    label: string
    options: LocationOption[]
  }

  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout>();
  const [locationSelectFilterValue, setLocationSelectFilterValue] = useState<string>("")
  const [locationSelectIsLoading, setLocationSelectIsLoading] = useState<boolean>(false)
  const [pinmapLocations, setPinmapLocations] = useState<LocationOption[]>([])
  const [pinmanLocations, setPinmanLocations] = useState<LocationOption[]>([])
  const [locationOptions, setLocationOptions] = useState<LocationOptionGroup[]>([])
  const [locationValue, setLocationValue] = useState<LocationOption>()

  useEffect(() => {
    const locationOptions = []

    if (pinmanLocations.length > 0) {
      locationOptions.push({
        label: 'Known Locations',
        options: pinmanLocations
      })
    }

    if (pinmapLocations.length > 0) {
      locationOptions.push({
        label: 'Pinball Map Locations (Create New Location)',
        options: pinmapLocations.filter((location: LocationOption) => {
          return pinmanLocations.find((pinmanLocation: LocationOption) => {
            return pinmanLocation.pinballMapId === location.pinballMapId
          }) === undefined
        })
      })
    }

    setLocationOptions(locationOptions)
  }, [pinmanLocations, pinmapLocations]);

  function updatePinmanLocations() {
    pinmanApi.locationsApi().locationsGet().then((locations) => {
      setPinmanLocations(
        locations.data.locations.map((location) : LocationOption => ({
            label: location.name,
            value: location.id,
            type: "pinman",
            pinballMapId: location.pinball_map_id.toString()
          })
        ))
    }).catch((e: AxiosError) => {
      const err = pinmanApi.parseError(e)
      console.error(err)
      setAlert({
        status: "error",
        title: err.title,
        detail: err.detail
      })
    }).finally(() => {
      setLocationSelectIsLoading(false)
    })
  }

  function updatePinballMapLocations(nameFilter: string) {
    if (nameFilter.length < 4) {
      return
    }

    pinballMapApi.locationsGet(nameFilter).then((locations) => {
      setPinmapLocations(
        locations.map((location): LocationOption => ({
            label: location.name,
            value: location.id.toString(),
            type: "pinmap",
            pinballMapId: location.id.toString()
          })
        ))
    }).catch((e: AxiosError) => {
      const err = pinballMapApi.parseError(e)
      console.error(err)
      setAlert({
        status: "error",
        title: "Error loading locations",
        detail: err.errors
      })
    }).finally(() => {
      setLocationSelectIsLoading(false)
    })
  }

  // useEffect to listen for changes in the input value
  useEffect(() => {
    if (locationSelectFilterValue.length < 4) {
      setLocationSelectIsLoading(false)
    } else {
      setLocationSelectIsLoading(true)
    }
    // Clear the previous timeout on each input change
    if (timeoutId)
      clearTimeout(timeoutId);

    // Start a new timeout that will call the delayedFunction after 3 seconds of inactivity
    const newTimeoutId = setTimeout(() => {
      updatePinmanLocations()
      updatePinballMapLocations(locationSelectFilterValue)
    }, 5000);
    setTimeoutId(newTimeoutId);

    // Cleanup the timeout on unmount to avoid memory leaks
    return () => {
      clearTimeout(newTimeoutId);
    };
  }, [locationSelectFilterValue]);

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
          <Select isMulti={false} onChange={onLocationChange} isLoading={locationSelectIsLoading}
                  onInputChange={(newValue) => {
                    setLocationSelectFilterValue(newValue)
                  }} options={locationOptions}/>
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
