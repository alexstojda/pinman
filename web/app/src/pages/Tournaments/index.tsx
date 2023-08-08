import {
  Box,
  Button,
  Card,
  CardBody,
  CardHeader,
  Flex,
  Grid,
  GridItem,
  Heading,
  HStack,
  Spacer,
  Text
} from "@chakra-ui/react";
import Footer from "../../layouts/footer";
import {Api, Tournament, useAuth} from "../../api";
import {useNavigate} from "react-router-dom";
import Alert, {AlertData} from "../../components/Alert";
import {useEffect, useState} from "react";
import {AxiosError} from "axios";

const api: Api = new Api();

type TournamentListPageProps = {
  createFormOpen?: boolean
}

export default function TournamentsListPage(props: TournamentListPageProps) {
  const auth = useAuth({});
  const navigate = useNavigate();
  const [tournaments, setTournaments] = useState<Tournament[]>([])
  const [alert, setAlert] = useState<AlertData>()

  useEffect(() => {
    api.tournamentsApi().tournamentsGet().then((getLeaguesResponse) => {
      setTournaments(getLeaguesResponse.data.tournaments)
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
  }, [props.createFormOpen])

  function renderList() {
    if (alert) {
      return (
        <Alert alert={alert}/>
      )
    }

    if (!tournaments || tournaments.length === 0) {
      return (
        <p>No leagues found, try creating one!</p>
      )
    }
    return (
      <Grid gap={4}>
        {tournaments.length > 0 && tournaments.map((tournament) => (
          <GridItem key={tournament.id}>
            <Card data-testid="tournament-card" direction={'row'}>
              <CardHeader display="flex" alignItems="center" justifyContent="flex-end">
                <Flex flexDirection="column" textAlign={'right'}>
                  <Heading size={'md'}>{tournament.name}</Heading>
                  {tournament.league &&
                    <Text fontSize={'0.8em'}>{tournament.league.name}</Text>
                  }
                </Flex>
              </CardHeader>
              {tournament.location &&
                <CardBody display="flex" alignItems="center" justifyContent="flex-end">
                  <Flex flexDirection="column" textAlign={'right'}>
                    <Text fontSize={'0.8em'}>{tournament.location.name}</Text>
                    <Text fontSize={'0.65em'}>{tournament.location.address}</Text>
                  </Flex>
                </CardBody>
              }
            </Card>
          </GridItem>
        ))}
      </Grid>
    )
  }

  return (
    <Box fontSize="xl" w={['90%', '85%', '80%']} maxW={800} mx="auto">
      <Box pb={10}>
        <HStack py={5}>
          <Heading
            textAlign={'center'}
            fontSize={'4xl'}
            fontWeight={'bold'}>
            Leagues
          </Heading>
          <Spacer/>
          {auth?.user &&
            <Button data-testid={'create-league'} size={'sm'} onClick={() => {
              navigate('/leagues/create')
            }}>
              Create League
            </Button>
          }
        </HStack>
        {renderList()}
        <Footer/>
      </Box>
    </Box>
  )
}
