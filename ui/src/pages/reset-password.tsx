import { useCallback, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import {
  Field,
  FieldAttributes,
  FieldProps,
  Form,
  Formik,
  FormikHelpers,
} from 'formik'
import * as Yup from 'yup'
import {
  Button,
  FormControl,
  FormErrorMessage,
  HStack,
  Input,
  Link as ChakraLink,
  Text,
  useToast,
  VStack,
} from '@chakra-ui/react'
import { Helmet } from 'react-helmet-async'
import AccountAPI from '@/api/account'
import Logo from '@/components/common/logo'
import DrawerLayout from '@/components/layout/drawer'
import variables from '@/theme/variables'

type FormValues = {
  newPassword: string
  newPasswordConfirmation: string
}

const ResetPasswordPage = () => {
  const params = useParams()
  const toast = useToast()
  const [isFormDisabled, setFormDisabled] = useState(false)
  const token = params.token as string
  const formSchema = Yup.object().shape({
    newPassword: Yup.string()
      .required('Password is required')
      .matches(
        /^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])(?=.{8,})/,
        'Must contain at least 8 characters, one Uppercase, one Lowercase, one number and one special character'
      ),
    newPasswordConfirmation: Yup.string()
      .oneOf([Yup.ref('newPassword'), undefined], 'Passwords do not match')
      .required('Confirm your password'),
  })

  const handleSubmit = useCallback(
    async (
      { newPassword }: FormValues,
      { setSubmitting }: FormikHelpers<FormValues>
    ) => {
      try {
        await AccountAPI.resetPassword({
          newPassword,
          token: token,
        })
        setFormDisabled(true)
        toast({
          title: 'Password updated successfully',
          status: 'success',
          isClosable: true,
        })
      } finally {
        setSubmitting(false)
      }
    },
    [token, toast]
  )

  return (
    <DrawerLayout>
      <>
        <Helmet>
          <title>Reset Password</title>
        </Helmet>
        <VStack spacing="25px" w="100%">
          <Logo className="w-16" isGlossy={true} />
          <h1 className="font-display text-2xl font-medium">Reset Password</h1>
          <Formik
            initialValues={{
              newPassword: '',
              newPasswordConfirmation: '',
            }}
            validationSchema={formSchema}
            validateOnBlur={false}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, isSubmitting }) => (
              <Form className="w-full">
                <VStack spacing={variables.spacing}>
                  <Field name="newPassword">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.newPassword && touched.newPassword
                            ? true
                            : false
                        }
                      >
                        <Input
                          {...field}
                          id="newPassword"
                          placeholder="New password"
                          type="password"
                          disabled={isSubmitting || isFormDisabled}
                        />
                        <FormErrorMessage>
                          {errors.newPassword}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Field name="newPasswordConfirmation">
                    {({ field }: FieldAttributes<FieldProps>) => (
                      <FormControl
                        isInvalid={
                          errors.newPasswordConfirmation &&
                          touched.newPasswordConfirmation
                            ? true
                            : false
                        }
                      >
                        <Input
                          {...field}
                          id="newPasswordConfirmation"
                          placeholder="Confirm new password"
                          type="password"
                          disabled={isSubmitting || isFormDisabled}
                        />
                        <FormErrorMessage>
                          {errors.newPasswordConfirmation}
                        </FormErrorMessage>
                      </FormControl>
                    )}
                  </Field>
                  <Button
                    variant="solid"
                    colorScheme="blue"
                    w="100%"
                    type="submit"
                    isLoading={isSubmitting}
                    disabled={isFormDisabled}
                  >
                    Reset password
                  </Button>
                </VStack>
              </Form>
            )}
          </Formik>
          <HStack spacing={variables.spacingXs}>
            <Text>Password already reset?</Text>
            <ChakraLink as={Link} to="/sign-in">
              Sign in
            </ChakraLink>
          </HStack>
        </VStack>
      </>
    </DrawerLayout>
  )
}

export default ResetPasswordPage
