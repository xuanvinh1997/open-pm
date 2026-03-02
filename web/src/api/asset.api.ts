import apiClient from './client'
import type { FileAsset } from '@/types/asset.types'

const base = (slug: string) => `/api/v1/workspaces/${slug}`

export const assetApi = {
  upload(slug: string, file: File, entityType: string, entityId: string) {
    const formData = new FormData()
    formData.append('file', file)
    formData.append('entity_type', entityType)
    formData.append('entity_id', entityId)
    return apiClient.post<FileAsset>(`${base(slug)}/assets`, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  list(slug: string, entityType: string, entityId: string) {
    return apiClient.get<{ results: FileAsset[] }>(
      `${base(slug)}/assets?entity_type=${entityType}&entity_id=${entityId}`,
    )
  },

  delete(slug: string, assetId: string) {
    return apiClient.delete(`${base(slug)}/assets/${assetId}`)
  },
}
