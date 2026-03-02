<script setup lang="ts">
import { watch, onBeforeUnmount } from 'vue'
import { useEditor, EditorContent, type JSONContent } from '@tiptap/vue-3'
import StarterKit from '@tiptap/starter-kit'
import LinkExtension from '@tiptap/extension-link'
import Placeholder from '@tiptap/extension-placeholder'
import UnderlineExtension from '@tiptap/extension-underline'
import TaskList from '@tiptap/extension-task-list'
import TaskItem from '@tiptap/extension-task-item'
import EditorToolbar from './EditorToolbar.vue'
import './editor.css'

interface Props {
  modelValue?: string
  json?: Record<string, unknown>
  placeholder?: string
  editable?: boolean
  toolbar?: 'full' | 'compact' | 'none'
  autofocus?: boolean
  minHeight?: string
  maxHeight?: string
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  placeholder: '',
  editable: true,
  toolbar: 'full',
  autofocus: false,
  minHeight: '100px',
})

const emit = defineEmits<{
  'update:modelValue': [html: string]
  'update:json': [json: JSONContent]
  'update:stripped': [text: string]
  blur: []
  focus: []
  submit: []
}>()

const editor = useEditor({
  content: props.modelValue || '',
  editable: props.editable,
  autofocus: props.autofocus,
  extensions: [
    StarterKit.configure({
      heading: { levels: [1, 2, 3] },
    }),
    LinkExtension.configure({
      openOnClick: false,
      HTMLAttributes: { rel: 'noopener noreferrer', target: '_blank' },
    }),
    Placeholder.configure({ placeholder: props.placeholder }),
    UnderlineExtension,
    TaskList,
    TaskItem.configure({ nested: true }),
  ],
  editorProps: {
    attributes: {
      class: 'prose prose-sm max-w-none focus:outline-none',
    },
    handleKeyDown: (_view, event) => {
      if ((event.metaKey || event.ctrlKey) && event.key === 'Enter') {
        emit('submit')
        return true
      }
      return false
    },
  },
  onUpdate: ({ editor: e }) => {
    const html = e.getHTML()
    emit('update:modelValue', html)
    emit('update:json', e.getJSON())
    emit('update:stripped', e.getText())
  },
  onFocus: () => emit('focus'),
  onBlur: () => emit('blur'),
})

// Sync external modelValue changes into the editor
watch(
  () => props.modelValue,
  (newVal) => {
    if (!editor.value) return
    const currentHtml = editor.value.getHTML()
    if (newVal !== currentHtml) {
      editor.value.commands.setContent(newVal || '', { emitUpdate: false })
    }
  },
)

watch(
  () => props.editable,
  (val) => {
    editor.value?.setEditable(val)
  },
)

onBeforeUnmount(() => {
  editor.value?.destroy()
})

defineExpose({ editor })
</script>

<template>
  <div
    v-if="editor"
    class="tiptap-editor rounded-md border border-custom-border-200 bg-custom-background-100 text-custom-text-200 overflow-hidden"
  >
    <EditorToolbar
      v-if="toolbar !== 'none' && editable"
      :editor="editor"
      :variant="toolbar === 'compact' ? 'compact' : 'full'"
    />
    <div
      :style="{
        '--editor-min-height': minHeight,
        maxHeight: maxHeight,
        overflowY: maxHeight ? 'auto' : undefined,
      }"
    >
      <EditorContent :editor="editor" />
    </div>
  </div>
</template>
