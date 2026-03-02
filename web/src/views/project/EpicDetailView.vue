<script setup lang="ts">
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// Epic detail now redirects to the issue detail page since epics are just issues with issue_type='epic'
const route = useRoute()
const router = useRouter()

const slug = route.params.workspaceSlug as string
const projectId = route.params.projectId as string
const epicId = route.params.epicId as string

onMounted(() => {
  router.replace(`/${slug}/projects/${projectId}/issues/${epicId}`)
})
</script>

<template>
  <div />
</template>
