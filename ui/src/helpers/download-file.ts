import { File } from '@/api/file'
import { getAccessTokenOrRedirect } from '@/infra/token'

export default async function downloadFile(file: File) {
  if (!file.original || file.type !== 'file') {
    return
  }
  const a: HTMLAnchorElement = document.createElement('a')
  a.href = `/proxy/api/v1/files/${file.id}/original${
    file.original.extension
  }?${new URLSearchParams({
    access_token: getAccessTokenOrRedirect(),
    download: 'true',
  })}`
  a.download = file.name
  a.style.display = 'none'
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
}
