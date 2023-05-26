import { IncomingMessage, ServerResponse } from 'http'
import path from 'path'
import { API_URL } from '@/config/config'
import { getTargetPath } from '@/infra/path'
import { File } from '@/api/file'
import { Token } from '@/api/token'

/*
  This method copies a resource from a source URL to a destination URL.

  Example implementation:

  - Extract the source and destination paths from the headers or request body.
  - Use fs.copyFile() to copy the file from the source to the destination.
  - Set the response status code to 204 if successful or an appropriate error code if the source file is not found or encountered an error.
  - Return the response.
 */
async function handleCopy(
  req: IncomingMessage,
  res: ServerResponse,
  token: Token
) {
  try {
    const sourceResult = await fetch(
      `${API_URL}/v1/files/get?path=${req.url}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
          'Content-Type': 'application/json',
        },
      }
    )
    const sourceFile: File = await sourceResult.json()

    const targetResult = await fetch(
      `${API_URL}/v1/files/get?path=${path.dirname(getTargetPath(req))}`,
      {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
          'Content-Type': 'application/json',
        },
      }
    )
    const targetFile: File = await targetResult.json()

    if (sourceFile.workspaceId !== targetFile.workspaceId) {
      res.statusCode = 400
      res.end()
      return
    }

    const copyResponse = await fetch(
      `${API_URL}/v1/files/${targetFile.id}/copy`,
      {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token.access_token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          ids: [sourceFile.id],
        }),
      }
    )
    const clones: File = await copyResponse.json()

    await fetch(`${API_URL}/v1/files/${clones[0].id}/rename`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token.access_token}`,
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        name: path.basename(getTargetPath(req)),
      }),
    })

    res.statusCode = 204
    res.end()
  } catch (err) {
    console.error(err)
    res.statusCode = 500
    res.end()
  }
}

export default handleCopy