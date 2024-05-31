import { useCallback, useMemo, useState } from 'react'
import { Button, ModalBody, ModalFooter } from '@chakra-ui/react'
import { OptionBase, Select, SingleValue } from 'chakra-react-select'
import cx from 'classnames'
import FileAPI from '@/client/api/file'
import InsightsAPI, { Language } from '@/client/api/insights'
import { swrConfig } from '@/client/options'
import { useAppDispatch, useAppSelector } from '@/store/hook'
import {
  creatingDidStart,
  creatingDidStop,
  modalDidClose,
} from '@/store/ui/insights'
import { reactSelectStyles } from '@/styles/react-select'

interface LanguageOption extends OptionBase {
  label: string
  value: string
}

const InsightsCreate = () => {
  const dispatch = useAppDispatch()
  const id = useAppSelector((state) =>
    state.ui.files.selection.length > 0
      ? state.ui.files.selection[0]
      : undefined,
  )
  const mutateMetadata = useAppSelector(
    (state) => state.ui.insights.mutateMetadata,
  )
  const mutateList = useAppSelector((state) => state.ui.files.mutate)
  const isCreating = useAppSelector((state) => state.ui.insights.isCreating)
  const [language, setLanguage] = useState<Language>()
  const { data: languages } = InsightsAPI.useGetLanguages(swrConfig())
  const { data: file } = FileAPI.useGet(id, swrConfig())
  const existingLanguage = useMemo<LanguageOption | undefined>(() => {
    if (file && languages && file.snapshot?.language) {
      const value = file.snapshot.language
        ? languages.filter((e) => e.id === file.snapshot?.language)[0]
        : undefined
      if (value) {
        setLanguage(value)
        return { value: value.id, label: value.name }
      }
    }
  }, [file, languages])

  const handleCreate = useCallback(async () => {
    if (id && language) {
      try {
        dispatch(creatingDidStart())
        await InsightsAPI.create(id, { languageId: language.id })
        mutateMetadata?.()
        mutateList?.()
      } catch (error) {
        dispatch(creatingDidStop())
      } finally {
        dispatch(creatingDidStop())
      }
    }
  }, [language, id, mutateMetadata, mutateList, dispatch])

  const handleLanguageChange = useCallback(
    (value: SingleValue<LanguageOption>) => {
      if (value?.value && languages) {
        setLanguage(languages.filter((e) => e.id === value.value)[0])
      }
    },
    [languages],
  )

  if (!id || !file || !languages) {
    return null
  }

  return (
    <>
      <ModalBody>
        <div
          className={cx(
            'flex',
            'flex-col',
            'items-center',
            'justify-center',
            'gap-1.5',
          )}
        >
          <p>
            Select the language to use for collecting insights.
            <br />
            During the process, text will be extracted using OCR (optical
            character recognition), and entities will be scanned using NER
            (named entity recognition).
          </p>
          {languages ? (
            <Select<LanguageOption, false>
              className={cx('w-full')}
              defaultValue={existingLanguage}
              options={languages.map((language) => ({
                value: language.id,
                label: language.name,
              }))}
              placeholder="Select Language"
              selectedOptionStyle="check"
              chakraStyles={reactSelectStyles()}
              isDisabled={isCreating}
              onChange={handleLanguageChange}
            />
          ) : null}
        </div>
      </ModalBody>
      <ModalFooter>
        <div className={cx('flex', 'flex-row', 'items-center', 'gap-1')}>
          <Button
            type="button"
            variant="outline"
            colorScheme="blue"
            isDisabled={isCreating}
            onClick={() => dispatch(modalDidClose())}
          >
            Cancel
          </Button>
          <Button
            type="button"
            variant="solid"
            colorScheme="blue"
            isLoading={isCreating}
            isDisabled={!language}
            onClick={handleCreate}
          >
            Collect Insights
          </Button>
        </div>
      </ModalFooter>
    </>
  )
}

export default InsightsCreate