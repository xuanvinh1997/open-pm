import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { pageApi } from '@/api/page.api'
import type { Page, CreatePageRequest, UpdatePageRequest } from '@/types/page.types'

export interface PageTreeNode extends Page {
  children: PageTreeNode[]
}

export const usePageStore = defineStore('page', () => {
  const pages = ref<Page[]>([])
  const currentPage = ref<Page | null>(null)
  const loading = ref(false)

  const treePages = computed<PageTreeNode[]>(() => {
    const map = new Map<string, PageTreeNode>()
    const roots: PageTreeNode[] = []

    // Initialize all nodes
    for (const page of pages.value) {
      map.set(page.id, { ...page, children: [] })
    }

    // Build tree
    for (const page of pages.value) {
      const node = map.get(page.id)!
      if (page.parent_id && map.has(page.parent_id)) {
        map.get(page.parent_id)!.children.push(node)
      } else {
        roots.push(node)
      }
    }

    return roots
  })

  async function fetchPages(slug: string, projectId: string) {
    loading.value = true
    try {
      const { data } = await pageApi.list(slug, projectId)
      pages.value = data.results
    } finally {
      loading.value = false
    }
  }

  async function fetchPage(slug: string, projectId: string, pageId: string) {
    const { data } = await pageApi.get(slug, projectId, pageId)
    currentPage.value = data
    return data
  }

  async function createPage(slug: string, projectId: string, data: CreatePageRequest) {
    const { data: page } = await pageApi.create(slug, projectId, data)
    pages.value.push(page)
    return page
  }

  async function updatePage(slug: string, projectId: string, pageId: string, data: UpdatePageRequest) {
    const { data: updated } = await pageApi.update(slug, projectId, pageId, data)
    const idx = pages.value.findIndex((p) => p.id === pageId)
    if (idx >= 0) pages.value[idx] = updated
    if (currentPage.value?.id === pageId) currentPage.value = updated
    return updated
  }

  async function deletePage(slug: string, projectId: string, pageId: string) {
    await pageApi.delete(slug, projectId, pageId)
    pages.value = pages.value.filter((p) => p.id !== pageId)
    if (currentPage.value?.id === pageId) currentPage.value = null
  }

  return {
    pages,
    currentPage,
    treePages,
    loading,
    fetchPages,
    fetchPage,
    createPage,
    updatePage,
    deletePage,
  }
})
