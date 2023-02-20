import { useCallback, useState } from 'react'
import {
  Button,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
  Text,
} from '@chakra-ui/react'
import GroupAPI, { Group } from '@/api/group'
import { User } from '@/api/user'
import variables from '@/theme/variables'

type GroupRemoveMemberProps = {
  group: Group
  user: User
  isOpen: boolean
  onClose?: () => void
  onCompleted?: () => void
}

const GroupRemoveMember = ({
  group,
  user,
  isOpen,
  onCompleted,
  onClose,
}: GroupRemoveMemberProps) => {
  const [loading, setLoading] = useState(false)

  const handleRemoveMember = useCallback(async () => {
    try {
      setLoading(true)
      await GroupAPI.removeMember(group.id, {
        userId: user.id,
      })
      onCompleted?.()
      onClose?.()
    } finally {
      setLoading(false)
    }
  }, [group, user, onCompleted, onClose])

  return (
    <Modal isOpen={isOpen} onClose={() => onClose?.()}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Remove Member</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <Text>
            Are you sure you would like to remove member{' '}
            <Text as="span" fontWeight="bold" whiteSpace="nowrap">
              {user.fullName}
            </Text>{' '}
            from group{' '}
            <Text as="span" fontWeight="bold" whiteSpace="nowrap">
              {group.name}
            </Text>
            ?
          </Text>
        </ModalBody>
        <ModalFooter>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            mr={variables.spacingSm}
            disabled={loading}
            onClick={() => onClose?.()}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="solid"
            colorScheme="red"
            isLoading={loading}
            onClick={handleRemoveMember}
          >
            Remove
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  )
}

export default GroupRemoveMember
