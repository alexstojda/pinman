import React, {useEffect, useState} from "react";
import {AxiosError} from "axios";
import {Api, pinballMap} from "../../api";
import {ActionMeta, Select, SingleValue} from "chakra-react-select";

export type LocationOption = {
  label: string
  value: string
  type: "pinman" | "pinmap"
  pinballMapId: string
}

type LocationOptionGroup = {
  label: string
  options: LocationOption[]
}

const pinmanApi = new Api()
const pinballMapApi = new pinballMap.Api()

type LocationSelectProps = {
  onChange?: (value: LocationOption) => void
  onError?: (error: any) => void
  value?: LocationOption
}

export function LocationSelect(props: LocationSelectProps) {
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout>();
  const [locationSelectFilterValue, setLocationSelectFilterValue] = useState<string>("")
  const [locationSelectIsLoading, setLocationSelectIsLoading] = useState<boolean>(false)
  const [pinmapLocations, setPinmapLocations] = useState<LocationOption[]>([])
  const [pinmanLocations, setPinmanLocations] = useState<LocationOption[]>([])
  const [locationOptions, setLocationOptions] = useState<LocationOptionGroup[]>([])

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
        locations.data.locations.map((location): LocationOption => ({
            label: location.name,
            value: location.id,
            type: "pinman",
            pinballMapId: location.pinball_map_id.toString()
          })
        ))
    }).catch((e: AxiosError) => {
      const err = pinmanApi.parseError(e)
      console.error(err)
      if (props.onError) {
        props.onError(err)
      }
    }).finally(() => {
      setLocationSelectIsLoading(false)
    })
  }

  function updatePinballMapLocations(nameFilter: string) {
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
      if (props.onError) {
        props.onError({
          detail: err.errors,
          title: "Pinball Map",
        })
      }
    }).finally(() => {
      setLocationSelectIsLoading(false)
    })
  }

  useEffect(() => {
    if (locationSelectFilterValue.length < 4) {
      setLocationSelectIsLoading(false)
      return
    } else {
      setLocationSelectIsLoading(true)
    }

    if (timeoutId)
      clearTimeout(timeoutId);

    const newTimeoutId = setTimeout(() => {
      updatePinmanLocations()
      updatePinballMapLocations(locationSelectFilterValue)
    }, 5000);
    setTimeoutId(newTimeoutId);

    return () => {
      clearTimeout(newTimeoutId);
    };
  }, [locationSelectFilterValue]);

  function onChange(newValue: SingleValue<LocationOption>, _: ActionMeta<LocationOption>) {
    if (props.onChange) {
      props.onChange(newValue as LocationOption)
    }
  }

  return (
    <Select
      isMulti={false}
      onChange={onChange}
      isLoading={locationSelectIsLoading}
      onInputChange={(newValue) => {
        setLocationSelectFilterValue(newValue)
      }}
      options={locationOptions}
      value={props.value}
    />
  )
}
