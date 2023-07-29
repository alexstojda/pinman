import {Api, League, useAuth} from "../../api";
import {useEffect, useState} from "react";
import {
  Button,
  Card,
  CardBody,
  CardHeader,
  Container,
  Grid,
  GridItem,
  Heading,
  HStack,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalHeader,
  ModalOverlay,
  Spacer,
} from "@chakra-ui/react";
import {useNavigate} from "react-router-dom";
import {AxiosError} from "axios";
import Alert, {AlertData} from "../../components/Alert";
import LeagueForm from "./Form";
import Footer from "../../layouts/footer";

const api: Api = new Api();

export type LeagueListPageProps = {
  createFormOpen?: boolean
}

export default function LeaguesPage(props: LeagueListPageProps) {
  const navigate = useNavigate();

  const auth = useAuth({requireAuth: props.createFormOpen});

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
  }, [props.createFormOpen])

  function renderList() {
    if (alert) {
      return (
        <Alert alert={alert}/>
      )
    }

    if (!leagues || leagues.length === 0) {
      return (
        <p>No leagues found, try creating one!</p>
      )
    }
    return (
      <Grid gap={4}>
        {leagues.length > 0 && leagues.map((league) => (
          <GridItem key={league.id}>
            <Card data-testid="league-card" direction={'row'}>
              <CardHeader>
                <Heading size={'md'}>{league.name}</Heading>
              </CardHeader>
              <CardBody>
                <p style={{textAlign: 'right'}}>{league.locationId}</p>
              </CardBody>
            </Card>
          </GridItem>
        ))}
      </Grid>
    )
  }

  return (
    <>
      <Modal isOpen={props.createFormOpen ?? false} onClose={() => navigate('/leagues')}>
        <ModalOverlay/>
        <ModalContent>
          <ModalHeader>Create new league</ModalHeader>
          <ModalCloseButton/>
          <ModalBody>
            <LeagueForm mode='create' onCancel={() => navigate('/leagues')}/>
          </ModalBody>
        </ModalContent>
      </Modal>

      <Container>
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
      </Container>
    </>
  )
}
