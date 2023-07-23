import {Api, League} from "../../api";
import {useEffect, useState} from "react";
import {Button, Card, CardBody, CardHeader, Container, Heading, HStack, Spacer} from "@chakra-ui/react";
import {useNavigate} from "react-router-dom";
import {AxiosError} from "axios";
import Alert, {AlertData} from "../../components/Alert";

const api: Api = new Api();

export default function LeagueListPage() {
  const navigate = useNavigate();

  const [leagues, setLeagues] = useState<League[]>([])
  const [alert, setAlert] = useState<AlertData>()

  useEffect(() => {
    api.leaguesApi().leaguesGet().then((getLeaguesResponse) => {
      setLeagues(getLeaguesResponse.data.leagues)
    }).catch((e: AxiosError) => {
      try {
        const err = api.parseError(e)
        console.error(err)
        setAlert({
          status: 'error',
          title: "Failed to list leagues",
          detail: err.detail
        })
      } catch (parseError) {
        setAlert({
          status: 'error',
          title: "Failed to list leagues",
          detail: 'Unknown error'
        })
      }
    })
  }, [])

  function renderList() {
    if (alert) {
      return (
        <Alert status={alert.status} title={alert.title} detail={alert.detail}/>
      )
    }

    if (leagues.length === 0) {
      return (
        <p>No leagues found, try creating one!</p>
      )
    }
    return (
      <>
        {leagues.length > 0 && leagues.map((league) => (
          <Card data-testid="league-card" key={league.id} direction={'row'}>
            <CardHeader>
              <Heading size={'md'}>{league.name}</Heading>
            </CardHeader>
            <CardBody>
              <p style={{textAlign: 'right'}}>{league.location}</p>
            </CardBody>
          </Card>
        ))}
      </>
    )
  }

  return (
    <Container>
      <HStack py={5}>
        <Heading
          textAlign={'center'}
          fontSize={'4xl'}
          fontWeight={'bold'}>
          Leagues
        </Heading>
        <Spacer/>
        <Button data-testid={'create-league'} size={'sm'} onClick={() => {navigate('/leagues/create')}}>
          Create League
        </Button>
      </HStack>
      {renderList()}
    </Container>
  )
}
