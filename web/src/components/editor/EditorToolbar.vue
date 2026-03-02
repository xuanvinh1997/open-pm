<script setup lang="ts">
import type { Editor } from '@tiptap/vue-3'
import {
  Bold,
  Italic,
  Underline,
  Strikethrough,
  Code,
  List,
  ListOrdered,
  Heading1,
  Heading2,
  Heading3,
  Quote,
  Minus,
  Link,
  ListChecks,
  Code2,
} from 'lucide-vue-next'

interface Props {
  editor: Editor
  variant?: 'full' | 'compact'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'full',
})

function setLink() {
  const previousUrl = props.editor.getAttributes('link').href
  const url = window.prompt('URL', previousUrl)
  if (url === null) return
  if (url === '') {
    props.editor.chain().focus().extendMarkRange('link').unsetLink().run()
    return
  }
  props.editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run()
}
</script>

<template>
  <div class="flex flex-wrap items-center gap-0.5 border-b border-custom-border-200 px-2 py-1">
    <!-- Headings (full only) -->
    <template v-if="variant === 'full'">
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('heading', { level: 1 }) }"
        title="Heading 1"
        @click="editor.chain().focus().toggleHeading({ level: 1 }).run()"
      >
        <Heading1 class="h-4 w-4" />
      </button>
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('heading', { level: 2 }) }"
        title="Heading 2"
        @click="editor.chain().focus().toggleHeading({ level: 2 }).run()"
      >
        <Heading2 class="h-4 w-4" />
      </button>
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('heading', { level: 3 }) }"
        title="Heading 3"
        @click="editor.chain().focus().toggleHeading({ level: 3 }).run()"
      >
        <Heading3 class="h-4 w-4" />
      </button>
      <div class="mx-1 h-4 w-px bg-custom-border-200" />
    </template>

    <!-- Inline formatting -->
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('bold') }"
      title="Bold"
      @click="editor.chain().focus().toggleBold().run()"
    >
      <Bold class="h-4 w-4" />
    </button>
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('italic') }"
      title="Italic"
      @click="editor.chain().focus().toggleItalic().run()"
    >
      <Italic class="h-4 w-4" />
    </button>
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('underline') }"
      title="Underline"
      @click="editor.chain().focus().toggleUnderline().run()"
    >
      <Underline class="h-4 w-4" />
    </button>
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('strike') }"
      title="Strikethrough"
      @click="editor.chain().focus().toggleStrike().run()"
    >
      <Strikethrough class="h-4 w-4" />
    </button>
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('code') }"
      title="Inline code"
      @click="editor.chain().focus().toggleCode().run()"
    >
      <Code class="h-4 w-4" />
    </button>

    <div class="mx-1 h-4 w-px bg-custom-border-200" />

    <!-- Lists -->
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('bulletList') }"
      title="Bullet list"
      @click="editor.chain().focus().toggleBulletList().run()"
    >
      <List class="h-4 w-4" />
    </button>
    <button
      v-if="variant === 'full'"
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('orderedList') }"
      title="Ordered list"
      @click="editor.chain().focus().toggleOrderedList().run()"
    >
      <ListOrdered class="h-4 w-4" />
    </button>
    <button
      v-if="variant === 'full'"
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('taskList') }"
      title="Task list"
      @click="editor.chain().focus().toggleTaskList().run()"
    >
      <ListChecks class="h-4 w-4" />
    </button>

    <!-- Full-only extras -->
    <template v-if="variant === 'full'">
      <div class="mx-1 h-4 w-px bg-custom-border-200" />
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('blockquote') }"
        title="Blockquote"
        @click="editor.chain().focus().toggleBlockquote().run()"
      >
        <Quote class="h-4 w-4" />
      </button>
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('codeBlock') }"
        title="Code block"
        @click="editor.chain().focus().toggleCodeBlock().run()"
      >
        <Code2 class="h-4 w-4" />
      </button>
      <button
        type="button"
        class="rounded p-1 hover:bg-custom-background-80 transition-colors"
        title="Horizontal rule"
        @click="editor.chain().focus().setHorizontalRule().run()"
      >
        <Minus class="h-4 w-4" />
      </button>
    </template>

    <div class="mx-1 h-4 w-px bg-custom-border-200" />

    <!-- Link -->
    <button
      type="button"
      class="rounded p-1 hover:bg-custom-background-80 transition-colors"
      :class="{ 'bg-custom-background-80 text-custom-text-100': editor.isActive('link') }"
      title="Link"
      @click="setLink"
    >
      <Link class="h-4 w-4" />
    </button>
  </div>
</template>
