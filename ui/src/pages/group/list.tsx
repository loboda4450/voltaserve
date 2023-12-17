import { useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import {
  Center,
  Heading,
  HStack,
  Link as ChakraLink,
  Stack,
  Table,
  Tbody,
  Td,
  Th,
  Thead,
  Tr,
  VStack,
  Text,
  Avatar,
  Badge,
} from '@chakra-ui/react'
import { SectionSpinner, variables } from '@koupr/ui'
import { Helmet } from 'react-helmet-async'
import GroupAPI, { SortOrder } from '@/client/api/group'
import { swrConfig } from '@/client/options'
import PagePagination from '@/components/common/page-pagination'
import { CreateGroupButton } from '@/components/top-bar/buttons'
import prettyDate from '@/helpers/pretty-date'
import { decodeQuery } from '@/helpers/query'
import usePagePagination from '@/hooks/page-pagination'

const GroupListPage = () => {
  const [searchParams] = useSearchParams()
  const query = decodeQuery(searchParams.get('q') as string)
  const { page, size, onPageChange, onSizeChange } = usePagePagination({
    localStoragePrefix: 'voltaserve',
    localStorageNamespace: 'group',
  })
  const {
    data: list,
    error,
    mutate,
  } = GroupAPI.useList(
    { query, page, size, sortOrder: SortOrder.Desc },
    swrConfig(),
  )

  useEffect(() => {
    mutate()
  }, [query, page, size, mutate])

  return (
    <>
      <Helmet>
        <title>Groups</title>
      </Helmet>
      <Stack
        direction="column"
        spacing={variables.spacing2Xl}
        pb={variables.spacing2Xl}
      >
        <Heading fontSize={variables.headingFontSize} pl={variables.spacingMd}>
          Groups
        </Heading>
        {error && (
          <Center h="300px">
            <Text>Failed to load groups.</Text>
          </Center>
        )}
        {!list && !error && <SectionSpinner />}
        {list && list.data.length === 0 && (
          <Center h="300px">
            <VStack spacing={variables.spacing}>
              <Text>There are no groups.</Text>
              <CreateGroupButton />
            </VStack>
          </Center>
        )}
        {list && list.data.length > 0 && (
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Name</Th>
                <Th>Organization</Th>
                <Th>Permission</Th>
                <Th>Date</Th>
              </Tr>
            </Thead>
            <Tbody>
              {list.data.map((g) => (
                <Tr key={g.id}>
                  <Td>
                    <HStack spacing={variables.spacing}>
                      <Avatar
                        name={g.name}
                        size="sm"
                        width="40px"
                        height="40px"
                      />
                      <ChakraLink
                        as={Link}
                        to={`/group/${g.id}/member`}
                        textDecoration="none"
                      >
                        {g.name}
                      </ChakraLink>
                    </HStack>
                  </Td>
                  <Td>
                    <ChakraLink
                      as={Link}
                      to={`/organization/${g.organization.id}/member`}
                      textDecoration="none"
                    >
                      {g.organization.name}
                    </ChakraLink>
                  </Td>
                  <Td>
                    <Badge>{g.permission}</Badge>
                  </Td>
                  <Td>{prettyDate(g.createTime)}</Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        )}
        {list && (
          <HStack alignSelf="end">
            <PagePagination
              totalPages={list.totalPages}
              page={page}
              size={size}
              onPageChange={onPageChange}
              onSizeChange={onSizeChange}
            />
          </HStack>
        )}
      </Stack>
    </>
  )
}

export default GroupListPage
