import { useCallback, useEffect, useState } from 'react'
import { useNavigate } from 'react-router-dom'
import {
  Button,
  FormControl,
  FormErrorMessage,
  Input,
  Modal,
  ModalBody,
  ModalCloseButton,
  ModalContent,
  ModalFooter,
  ModalHeader,
  ModalOverlay,
} from '@chakra-ui/react'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import cx from 'classnames'
import GroupAPI, { Group } from '@/client/api/group'

export type GroupDeleteProps = {
  open: boolean
  group: Group
  onClose?: () => void
}

type FormValues = {
  name: string
}

const GroupDelete = ({ open, group, onClose }: GroupDeleteProps) => {
  const navigate = useNavigate()
  const [isModalOpen, setIsModalOpen] = useState(false)
  const formSchema = Yup.object().shape({
    name: Yup.string()
      .required('Confirmation is required')
      .oneOf([group.name], 'Invalid group name'),
  })

  useEffect(() => {
    setIsModalOpen(open)
  }, [open])

  const handleSubmit = useCallback(
    async (_: FormValues, { setSubmitting }: FormikHelpers<FormValues>) => {
      setSubmitting(true)
      try {
        await GroupAPI.delete(group.id)
        navigate('/group')
        onClose?.()
      } finally {
        setSubmitting(false)
      }
    },
    [group.id, navigate, onClose],
  )

  return (
    <Modal
      isOpen={isModalOpen}
      onClose={() => onClose?.()}
      closeOnOverlayClick={false}
    >
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Delete Group</ModalHeader>
        <ModalCloseButton />
        <Formik
          initialValues={{ name: '' }}
          validationSchema={formSchema}
          validateOnBlur={false}
          onSubmit={handleSubmit}
        >
          {({ errors, touched, isSubmitting }) => (
            <Form>
              <ModalBody>
                <div className={cx('flex', 'flex-col', 'gap-1.5')}>
                  <span>Are you sure you would like to delete this group?</span>
                  <span>
                    Please type <b>{group.name}</b> to confirm.
                  </span>
                  <Field name="name">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={errors.name && touched.name ? true : false}
                      >
                        <Input {...field} disabled={isSubmitting} />
                        <FormErrorMessage>{errors.name}</FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                </div>
              </ModalBody>
              <ModalFooter>
                <div
                  className={cx('flex', 'flex-row', 'items-center', 'gap-1')}
                >
                  <Button
                    type="button"
                    variant="outline"
                    colorScheme="blue"
                    disabled={isSubmitting}
                    onClick={() => onClose?.()}
                  >
                    Cancel
                  </Button>
                  <Button
                    type="submit"
                    variant="solid"
                    colorScheme="red"
                    disabled={isSubmitting}
                    isLoading={isSubmitting}
                  >
                    Delete Permanently
                  </Button>
                </div>
              </ModalFooter>
            </Form>
          )}
        </Formik>
      </ModalContent>
    </Modal>
  )
}

export default GroupDelete